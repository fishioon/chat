package main

//import "github.com/golang/protobuf/proto"

import (
	"flag"
	"fmt"
	pb "github.com/fishioon/chat/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	//"google.golang.org/grpc/grpclog"
	//"time"
)

var (
	grpcAddr  = flag.String("grpcaddr", ":8001", "listen grpc addr")
	httpAddr  = flag.String("addr", ":8000", "listen http addr")
	debugAddr = flag.String("debugaddr", ":8002", "listen debug addr")
)

func listenGRPC(listenAddr string) error {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	cs := NewChatServer()
	pb.RegisterChatServer(grpcServer, cs)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("serveGRPC err:", err)
		}
	}()
	return nil
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	listenGRPC(*grpcAddr)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterChatHandlerFromEndpoint(ctx, mux, *grpcAddr, opts)
	if err != nil {
		log.Fatal(err)
	}
	go http.ListenAndServe(*debugAddr, nil)
	fmt.Println("listening")
	http.ListenAndServe(*httpAddr, WebsocketProxy(mux))
}
