package server

import (
	"backup-client/conf"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func RegistrarServer(c *conf.Bootstrap, logger log.Logger) (*etcd.Registry, func(), error) {
	// new etcd client
	logr := log.NewHelper(logger)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: c.Registry.GetEtcd(),
		Username:  c.Registry.GetEtcdUserName(),
		Password:  c.Registry.GetEtcdPassword(),
	})
	if err != nil {
		logr.Errorf("etcd registration error:[%v]", err)
		return nil, nil, err
	}
	logr.Info("etcd registration successful")
	// new reg with etcd client
	return etcd.New(client, etcd.MaxRetry(int(c.Registry.MaxRetry))), func() {
	}, nil
}
