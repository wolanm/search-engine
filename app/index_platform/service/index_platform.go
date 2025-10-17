package service

import (
	"context"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/kevwan/mapreduce/v2"
	"github.com/spf13/cast"
	"github.com/wolanm/search-engine/app/index_platform/analyzer"
	"github.com/wolanm/search-engine/app/index_platform/indexplatform_logger"
	"github.com/wolanm/search-engine/app/index_platform/input_data"
	"github.com/wolanm/search-engine/consts"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
	"github.com/wolanm/search-engine/types"
	"google.golang.org/grpc/metadata"
	"io"
	"sort"
	"strings"
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
	md, ok := metadata.FromIncomingContext(ctx)
	streamRespFunc := func(code int64, message string) {
		_ = stream.SendAndClose(&pb.UploadResponse{
			Code:    code,
			Message: message,
		})
	}
	if !ok {
		err = fmt.Errorf("metadata not provided")
		indexplatform_logger.Logger.Error("UploadFile metadata.FromIncomingContext failed: ", err)
		streamRespFunc(int64(consts.InvalidParam), "metadata not provided")
		return
	}
	indexplatform_logger.Logger.Info("Received file: : ", md["filename"])

	// 如果 reduce 会并发运行，则需要考虑使用 concurrent map，当前 mapreduce 的 reduce 是非并发的
	inverted_index := make(map[string]*roaring.Bitmap)
	// TODO: 使用 trie 树 dictTrie := trie.NewTrie()
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		fileContent := make([]byte, 0)
		for {
			fileChunk, err := stream.Recv()
			if err == io.EOF {
				break
			}

			fileContent = append(fileContent, fileChunk.Content...)
		}

		streamRespFunc(int64(consts.Success), "Upload File Success")
		// TODO: 生成 document ID
		source <- fileContent
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(err error)) {
		// TODO: 控制并发

		keyValueList := make([]*types.KeyValue, 0, consts.DEFAULT_KV_LIST_CAPACITY)
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
					// TODO: 构建前缀树
				}
			}

			// TODO: 构建正排索引
			go func(docStruct *types.Document) {}(docStruct)

			// shuffle 排序
			sort.Sort(types.ByKey(keyValueList))
			writer.Write(keyValueList)
		}

	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		// 构建倒排索引
		for kvList := range pipe {
			for _, kv := range kvList {
				var docIds *roaring.Bitmap
				if docIds, ok = inverted_index[kv.Key]; !ok {
					inverted_index[kv.Key] = roaring.New()
					docIds = inverted_index[kv.Key]
				}

				docIds.AddInt(cast.ToInt(kv.Key))
			}
		}
	})

	// TODO: 存储倒排索引
	go func() {
		// TODO: 实现链路追踪后，这里的 ctx 要 clone
		ctx := context.Background()
		err := storeInvertedIndex(ctx, inverted_index)
		if err != nil {
			indexplatform_logger.Logger.Error("storeInvertedIndex error: ", err)
		}
	}()
	// TODO: 存储前缀树

	return nil
}

func (s *IndexPlatformSrv) DownloadFile(file *pb.FileRequest, req pb.IndexPlatformService_DownloadFileServer) (err error) {
	return nil
}

func storeInvertedIndex(ctx context.Context, inverted_index map[string]*roaring.Bitmap) (err error) {
	
}
