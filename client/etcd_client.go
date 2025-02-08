package client

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"pod-service-relations/logging"
	"time"
)

func NewEtcdClient(ctx context.Context, etcdServerHost string, etcdServerPort int) (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", etcdServerHost, etcdServerPort)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("get etcd client got error: %s", err))
		return nil, err
	}
	return client, err
}
