package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	EtcdClient, err := clientv3.New(clientv3.Config{
		DialTimeout: 3 * time.Second,
		Endpoints:   []string{"192.168.11.207:2379"},
		Username:    "",
		Password:    "",
	})

	r := etcd.New(EtcdClient)

	watcher, err := r.Watch(context.TODO(), "local")
	if err != nil {
		panic(err)
	}
	for {

		instances, err := watcher.Next()
		fmt.Println(instances, err)
	}

}
