package discovery

import "time"

type ServiceNode struct {
	NodeID        string            `json:"node_id"`        // 唯一节点ID
	ServiceName   string            `json:"service_name"`   // 服务名称
	ServiceType   string            `json:"service_type"`   // 服务类型
	Endpoint      string            `json:"endpoint"`       // gRPC 地址
	Version       string            `json:"version"`        // 服务版本
	Metadata      map[string]string `json:"metadata"`       // 元数据
	Status        string            `json:"status"`         // 健康状态
	Load          int32             `json:"load"`           // 负载指标
	LastHeartbeat time.Time         `json:"last_heartbeat"` // 最后心跳
}
