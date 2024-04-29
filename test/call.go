package main

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"

	helloworld "backup-client/api/helloworld/v1"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	EtcdClient, err := clientv3.New(clientv3.Config{
		DialTimeout: 3 * time.Second,
		Endpoints:   []string{"192.168.11.207:2379"},
		Username:    "",
		Password:    "",
	})

	r := etcd.New(EtcdClient)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///local"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn)

	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///local"),
		http.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := helloworld.NewGreeterHTTPClient(hConn)

	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}
}

func callGRPC(client helloworld.GreeterClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(client helloworld.GreeterHTTPClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}
