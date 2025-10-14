package service

import (
	"context"
	"fmt"
	"github.com/wolanm/search-engine/app/index_platform/indexplatform_logger"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
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
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			break
		}

		indexplatform_logger.Logger.Info("recv file content: %s", fileChunk.String())
	}
	fmt.Println("UploadFile called")

	_ = stream.SendAndClose(&pb.UploadResponse{
		Code:    0,
		Message: "Upload File Success",
	})

	return nil
}

func (s *IndexPlatformSrv) DownloadFile(file *pb.FileRequest, req pb.IndexPlatformService_DownloadFileServer) (err error) {
	return nil
}
