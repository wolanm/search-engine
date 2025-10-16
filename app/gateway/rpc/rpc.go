package rpc

import (
	"context"
	"github.com/wolanm/search-engine/app/gateway/gateway_logger"
	"github.com/wolanm/search-engine/config"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/grpc_client"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
	"github.com/wolanm/search-engine/util/discovery"
	"io"
	"mime/multipart"
)

var indexPlatformCli *grpc_client.IndexPlatFormClient

func Init() {
	// 服务发现
	err := discovery.RegisterResolver([]string{config.Conf.Etcd.Address}, gateway_logger.Logger)
	if err != nil {
		panic(err)
	}

	indexPlatformCli, err = grpc_client.NewIndexPlatformClient(config.Conf.Services[consts.IndexPlatform].Name)
	if err != nil {
		panic(err)
	}
}

func UploadFile(ctx context.Context, file multipart.File, filesize int64) (resp *pb.UploadResponse, err error) {
	stream, err := indexPlatformCli.Client.UploadFile(ctx)

	// 读取文件，通过 stream 传输
	if err != nil {
		gateway_logger.Logger.Error("indexPlatformCli.Client.UploadFile failed: ", err)
		return
	}

	buf := make([]byte, 1024*1024) // 1MB chunks
	for {
		n, errx := file.Read(buf)
		if errx == io.EOF {
			break
		}

		if err = stream.Send(&pb.FileChunk{Content: buf[:n]}); err != nil {
			gateway_logger.Logger.Error("stream.Send failed: ", err)
			return
		}
	}

	resp, err = stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return resp, nil
}
