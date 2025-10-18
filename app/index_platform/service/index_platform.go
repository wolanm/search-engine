package service

import (
	"context"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/cespare/xxhash"
	"github.com/kevwan/mapreduce/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/wolanm/search-engine/app/index_platform/analyzer"
	"github.com/wolanm/search-engine/app/index_platform/indexplatform_logger"
	"github.com/wolanm/search-engine/app/index_platform/input_data"
	"github.com/wolanm/search-engine/app/index_platform/kfk"
	"github.com/wolanm/search-engine/consts"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
	"github.com/wolanm/search-engine/repository/inverted_db"
	"github.com/wolanm/search-engine/repository/trie_db"
	"github.com/wolanm/search-engine/types"
	"github.com/wolanm/search-engine/util/path"
	"github.com/wolanm/search-engine/util/trie"
	"google.golang.org/grpc/metadata"
	"io"
	"sort"
	"strings"
	"sync"
)

type IndexPlatformSrv struct {
	pb.UnimplementedIndexPlatformServiceServer
}

func NewIndexPlatformSrv() *IndexPlatformSrv {
	return &IndexPlatformSrv{}
}

func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	return nil, nil
}

func (s *IndexPlatformSrv) UploadFile(stream pb.IndexPlatformService_UploadFileServer) (err error) {
	ctx := stream.Context()
	streamRespFunc := func(code int64, message string) {
		_ = stream.SendAndClose(&pb.UploadResponse{
			Code:    code,
			Message: message,
		})
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = fmt.Errorf("metadata not provided")
		indexplatform_logger.Logger.Error("UploadFile metadata.FromIncomingContext failed: ", err)
		streamRespFunc(int64(consts.InvalidParam), "metadata not provided")
		return
	}
	indexplatform_logger.Logger.Info("Received file: ", md["filename"][0])

	// 接收和处理需要分离，否则 grpc 调用时间与 mapreduce 处理时间一样长
	fileContent := make([]byte, 0)
	for {
		fileChunk, errx := stream.Recv()
		if errx == io.EOF {
			break
		}

		fileContent = append(fileContent, fileChunk.Content...)
	}
	streamRespFunc(int64(consts.Success), "Upload File Success")
	indexplatform_logger.Logger.Infof("read %s finish", md["filename"])

	go buildIndex(md["filename"][0], fileContent)
	return nil
}

func (s *IndexPlatformSrv) DownloadFile(file *pb.FileRequest, req pb.IndexPlatformService_DownloadFileServer) (err error) {
	return nil
}

func buildIndex(filename string, fileContent []byte) {
	// 构建正排索引

	// 构建倒排索引, 向量索引
	kfk.DocDataToKfk()

	// 如果 reduce 会并发运行，则需要考虑使用 concurrent map，当前 mapreduce 的 reduce 是非并发的
	invertedIndex := make(map[string]*roaring.Bitmap)
	dictTrie := trie.NewTrie()
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		source <- fileContent
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(err error)) {
		// TODO: 控制并发
		ch := make(chan struct{}, consts.ConcurrentMapWorker)
		var wg sync.WaitGroup

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
				indexplatform_logger.Logger.Error("GseCutForBuildIndex error: ", err)
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

			ch <- struct{}{}
			wg.Add(1)

			// TODO: 构建正排索引，向量索引
			go func(docStruct *types.Document) {
				// kafka 发送数据
				defer wg.Done()
				<-ch
			}(docStruct)
		}
		wg.Wait()

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
			indexplatform_logger.Logger.Error("storeInvertedIndex error: ", err)
		} else {
			indexplatform_logger.Logger.Infof("storeInvertedIndex %s success", filename)
		}
	}()

	go func() {
		newCtx := context.Background()
		err := storeTrie(newCtx, dictTrie)
		if err != nil {
			indexplatform_logger.Logger.Error("storeTrie error: ", err)
		} else {
			indexplatform_logger.Logger.Infof("storeTrie %s success", filename)
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
			indexplatform_logger.Logger.Error("close db failed: ", err)
		}
	}()

	// 遍历 inverted_index, 存储到 db 中
	for word, bitmap := range invertedIndex {
		data, errx := bitmap.MarshalBinary()
		if errx != nil {
			indexplatform_logger.Logger.Error("marshal bitmap failed: ", errx)
			continue
		}
		errx = invertedDB.StorageInverted(word, data)
		if errx != nil {
			indexplatform_logger.Logger.Error("StorageInverted failed: ", errx)
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
			indexplatform_logger.Logger.Error("Close TrieDB failed: ", errx)
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
