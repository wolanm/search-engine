package discovery

import (
	"context"
	"errors"
	log "github.com/wolanm/search-engine/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

// CustomFormatter 自定义日志格式器
type CustomFormatter struct {
	Module string // 模块名称
}

type ResolverBuilder struct {
	EtcdAddrs []string
	schema    string
	logger    *log.Logger
}

type Resolver struct {
	closeCh chan struct{}
	cc      resolver.ClientConn
	cli     *clientv3.Client
	key     string // 解析的地址
	watchCh clientv3.WatchChan
	addrMap map[string]resolver.Address // 服务提供者存储到 etcd 的 <key, value> 对
	logger  *log.Logger
}

func (rb *ResolverBuilder) Scheme() string {
	return rb.schema
}

func (rb *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		closeCh: make(chan struct{}),
		key:     "/" + target.Endpoint(),
		cc:      cc,
		logger:  rb.logger,
	}
	rb.logger.Infof("target: %s", target.String())

	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   rb.EtcdAddrs,
		DialTimeout: time.Duration(5) * time.Second,
	})

	if err != nil {
		rb.logger.Error("new etcd client error: ", err)
		return nil, err
	}

	if err := r.start(); err != nil {
		rb.logger.Error("resolver start failed, error: ", err)
		return nil, err
	}

	rb.logger.Info("resolver start success")
	return r, nil
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
	if err := r.resolve(); err != nil {
		r.logger.Error("Error in ResolveNow: ", err)
	}
}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// start 开始解析并监听 key
func (r *Resolver) start() error {
	var err error
	if err = r.resolve(); err != nil {
		return err
	}

	go r.watcher()

	return err
}

// resolve 执行一次解析
func (r *Resolver) resolve() error {
	resp, err := r.cli.Get(context.Background(), r.key, clientv3.WithPrefix())
	if err != nil {
		r.logger.Errorf("get %s --prefix failed: %v", r.key, err)
		return err
	}

	// 清空地址
	r.addrMap = map[string]resolver.Address{}
	for k, v := range resp.Kvs {
		// 解析得到 node
		node, err := ParseValue(v.Value)
		if err != nil {
			r.logger.Warnf("parse %s failed", string(v.Value))
			continue
		}

		r.addrMap[string(k)] = resolver.Address{Addr: node.Endpoint}
	}

	err = r.cc.UpdateState(resolver.State{Addresses: ConvertToGRPCAddress(r.addrMap)})
	if err != nil {
		r.logger.Error("update state failed: ", err)
	}
	return err
}

func (r *Resolver) watcher() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.cli.Watch(context.Background(), r.key, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.resolve(); err != nil {
				r.logger.Error("resolve  failed", err)
			}
		}
	}
}

func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			// 修改 addr
			node, err := ParseValue(ev.Kv.Value)
			if err != nil {
				r.logger.Warnf("parse %s failed", string(ev.Kv.Value))
				continue
			}

			r.addrMap[string(ev.Kv.Key)] = resolver.Address{Addr: node.Endpoint}
			_ = r.cc.UpdateState(resolver.State{Addresses: ConvertToGRPCAddress(r.addrMap)})
		case clientv3.EventTypeDelete:
			delete(r.addrMap, string(ev.Kv.Key))
			_ = r.cc.UpdateState(resolver.State{Addresses: ConvertToGRPCAddress(r.addrMap)})
		}
	}
}

func RegisterResolver(etcdAddrs []string, logger *log.Logger) error {
	if nil == logger {
		return errors.New("please pass a valid logger")
	}

	resolver.Register(&ResolverBuilder{EtcdAddrs: etcdAddrs, logger: logger, schema: "etcd"})
	return nil
}
