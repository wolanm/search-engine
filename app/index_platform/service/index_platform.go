package service

import (
	"context"
	"fmt"
	"github.com/wolanm/search-engine/app/index_platform/indexplatform_logger"
	"github.com/wolanm/search-engine/app/index_platform/kfk"
	"github.com/wolanm/search-engine/consts"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
	"github.com/wolanm/search-engine/types"
	"google.golang.org/grpc/metadata"
	"io"
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
	// 构建正排索引，倒排索引, 向量索引
	fileInfo := &types.FileInfo{
		Filename: filename,
		Content:  fileContent,
	}
	err := kfk.DocDataToKfk(fileInfo)
	if err != nil {
		indexplatform_logger.Logger.Error("send document to kafka failed", err)
	}
	return
}
