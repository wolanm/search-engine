package service

import (
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
)

type IndexPlatformSrv struct {
	pb.UnimplementedIndexPlatformServiceServer
}

//func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
//
//}

func (s *IndexPlatformSrv) UploadFile(stream pb.IndexPlatformService_UploadFileServer) (err error) {
	return nil
}

func (s *IndexPlatformSrv) DownloadFile(file *pb.FileRequest, req pb.IndexPlatformService_DownloadFileServer) (err error) {

	return nil
}
