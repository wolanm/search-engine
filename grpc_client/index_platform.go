package grpc_client

import (
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
)

type IndexPlatFormClient struct {
	Client pb.IndexPlatformServiceClient
}

func NewIndexPlatformClient(serviceName string) (*IndexPlatFormClient, error) {
	conn, err := connectServer(serviceName)

	return &IndexPlatFormClient{
		Client: pb.NewIndexPlatformServiceClient(conn),
	}, err
}
