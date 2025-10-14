package discovery

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/wolanm/search-engine/logger"
	"go.etcd.io/etcd/client/v3"
)

type ServiceRegistry struct {
	client        *clientv3.Client
	leaseID       clientv3.LeaseID
	node          *ServiceNode
	stopChan      chan struct{}
	keepaliveChan <-chan *clientv3.LeaseKeepAliveResponse
	ttl           int64
	logger        *log.Logger
}

// NewServiceRegistry /*
func NewServiceRegistry(endpoints []string, logger *log.Logger) (*ServiceRegistry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}
	return &ServiceRegistry{
		client:   client,
		stopChan: make(chan struct{}),
		logger:   logger,
	}, nil
}

// RegisterService 服务注册
func (sr *ServiceRegistry) RegisterService(node *ServiceNode, ttl int64) error {
	sr.node = node
	sr.ttl = ttl

	if err := sr.register(); err != nil {
		return err
	}

	go sr.keepalive()

	return nil
}

func (sr *ServiceRegistry) register() error {
	// 创建租约
	resp, err := sr.client.Grant(context.Background(), sr.ttl)
	if err != nil {
		return err
	}

	sr.leaseID = resp.ID

	// 开启自动续约
	if sr.keepaliveChan, err = sr.client.KeepAlive(context.Background(), sr.leaseID); err != nil {
		return err
	}

	// 序列化节点信息，作为服务的 value
	nodeData, err := json.Marshal(sr.node)
	if err != nil {
		return err
	}

	// 设置服务信息
	serviceKey := sr.node.BuildRegistryPath()
	if _, err = sr.client.Put(context.Background(), serviceKey, string(nodeData), clientv3.WithLease(sr.leaseID)); err != nil {
		return err
	}

	return nil
}

func (sr *ServiceRegistry) stop() {
	sr.stopChan <- struct{}{}
}

func (sr *ServiceRegistry) unregister() error {
	if _, err := sr.client.Delete(context.Background(), sr.node.BuildRegistryPath()); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRegistry) keepalive() {
	ticker := time.NewTicker(time.Duration(sr.ttl) * time.Second)
	for {
		select {
		case <-sr.stopChan:
			if err := sr.unregister(); err != nil {
				sr.logger.Error("unregister failed, error: ", err)
			}

			if _, err := sr.client.Revoke(context.Background(), sr.leaseID); err != nil {
				sr.logger.Error("revoke failed, error: ", err)
			}

			return

		case res := <-sr.keepaliveChan:
			if nil == res {
				if err := sr.register(); err != nil {
					sr.logger.Error("register failed, error: ", err)
				}
			}
		case <-ticker.C:
			if sr.keepaliveChan == nil {
				if err := sr.register(); err != nil {
					sr.logger.Error("register failed, error: ", err)
				}
			}
		}
	}
}
