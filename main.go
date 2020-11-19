package main

import (
	"github.com/xiiiew/websocket-push-go/common"
	"github.com/xiiiew/websocket-push-go/handles"
	pb "github.com/xiiiew/websocket-push-go/protos"
	"github.com/xiiiew/websocket-push-go/routes"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// http/ws服务
	router := routes.Register()
	go http.ListenAndServe(":"+common.Config.HttpWsPort, router)

	// rpc服务
	lis, err := net.Listen("tcp", ":"+common.Config.RpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterWsPushServer(grpcServer, &handles.WsPushServer{})
	grpcServer.Serve(lis)

	select {}
}
