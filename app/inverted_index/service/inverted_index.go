package service

import (
	"context"
	"github.com/RoaringBitmap/roaring"
	"github.com/cespare/xxhash"
	"github.com/kevwan/mapreduce/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/wolanm/search-engine/app/index_platform/input_data"
	"github.com/wolanm/search-engine/app/inverted_index/analyzer"
	"github.com/wolanm/search-engine/app/inverted_index/inverted_index_logger"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/repository/inverted_db"
	"github.com/wolanm/search-engine/repository/trie_db"
	"github.com/wolanm/search-engine/types"
	"github.com/wolanm/search-engine/util/path"
	"github.com/wolanm/search-engine/util/trie"
	"sort"
	"strings"
)

func BuildInvertedIndex(fileInfo *types.FileInfo) {
	inverted_index_logger.Logger.Info("build inverted index: ", fileInfo.Filename)

	// 如果 reduce 会并发运行，则需要考虑使用 concurrent map，当前 mapreduce 的 reduce 是非并发的
	invertedIndex := make(map[string]*roaring.Bitmap)
	dictTrie := trie.NewTrie()
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		source <- fileInfo.Content
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(err error)) {
		// TODO: 计算 tf-idf
		// TODO: 控制并发

		keyValueList := make([]*types.KeyValue, 0, consts.DefaultKvListCapacity)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines {
			// 分词
			docStruct, _ := input_data.Doc2Struct(line)
			if docStruct.DocId == 0 {
				continue
			}

			tokens, err := analyzer.GseCutForBuildIndex(docStruct.DocId, line)
			if err != nil {
				inverted_index_logger.Logger.Error("GseCutForBuildIndex error: ", err)
				continue
			} else {
				for _, v := range tokens {
					if v.Token == "" || v.Token == " " {
						continue
					}

					keyValueList = append(keyValueList, &types.KeyValue{Key: v.Token, Value: cast.ToString(docStruct.DocId)})
					dictTrie.Insert(v.Token)
				}
			}

		}

		// shuffle 排序
		sort.Sort(types.ByKey(keyValueList))
		writer.Write(keyValueList)
	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		// 构建倒排索引
		for kvList := range pipe {
			for _, kv := range kvList {
				var docIds *roaring.Bitmap
				var ok bool
				if docIds, ok = invertedIndex[kv.Key]; !ok {
					invertedIndex[kv.Key] = roaring.New()
					docIds = invertedIndex[kv.Key]
				}

				docIds.AddInt(cast.ToInt(kv.Key))
			}
		}
	})

	go func() {
		// TODO: 实现链路追踪后，这里的 ctx 要 clone
		newCtx := context.Background()
		err := storeInvertedIndex(newCtx, invertedIndex)
		if err != nil {
			inverted_index_logger.Logger.Error("storeInvertedIndex error: ", err)
		} else {
			inverted_index_logger.Logger.Infof("storeInvertedIndex %s success", fileInfo.Filename)
		}
	}()

	go func() {
		newCtx := context.Background()
		err := storeTrie(newCtx, dictTrie)
		if err != nil {
			inverted_index_logger.Logger.Error("storeTrie error: ", err)
		} else {
			inverted_index_logger.Logger.Infof("storeTrie %s success", fileInfo.Filename)
		}
	}()
}

func storeInvertedIndex(ctx context.Context, invertedIndex map[string]*roaring.Bitmap) (err error) {
	// 暂不考虑分片，内部知识库文档数据量比较小，且写入场景比较少，主要是读
	dbPath := path.GetInvertedDBPath()
	invertedDB, err := inverted_db.NewInvertedDB(dbPath)
	if err != nil {
		return err
	}
	defer func() {
		err = invertedDB.Close()
		if err != nil {
			inverted_index_logger.Logger.Error("close db failed: ", err)
		}
	}()

	// 遍历 inverted_index, 存储到 db 中
	for word, bitmap := range invertedIndex {
		data, errx := bitmap.MarshalBinary()
		if errx != nil {
			inverted_index_logger.Logger.Error("marshal bitmap failed: ", errx)
			continue
		}
		errx = invertedDB.StorageInverted(word, data)
		if errx != nil {
			inverted_index_logger.Logger.Error("StorageInverted failed: ", errx)
			continue
		}
	}

	return nil
}

func storeTrie(ctx context.Context, dict *trie.Trie) error {
	dbPath := path.GetTrieDBPath()
	trieDB, err := trie_db.NewTrieDB(dbPath)
	if err != nil {
		return err
	}

	defer func() {
		errx := trieDB.Close()
		if errx != nil {
			inverted_index_logger.Logger.Error("Close TrieDB failed: ", errx)
		}
	}()

	err = trieDB.StorageDict(dict)
	if err != nil {
		return errors.WithMessage(err, "storageDict failed")
	}

	return nil
}

// iHash 分片存储使用
func iHash(key string) uint64 { //  nolint:golint,unused
	hash := xxhash.Sum64String(key)
	return hash
}
