package grpc_client

import (
	"fmt"
	"github.com/wolanm/search-engine/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	addr := fmt.Sprintf("%s%s/%s", consts.ServiceDomain, consts.ServicePrefix, serviceName)

	conn, err = grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")),
	)
	return
}
