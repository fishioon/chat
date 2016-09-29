package main

import (
	"flag"
	pb "github.com/fishioon/chat/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	address := flag.String("address", "localhost:8001", "server address")
	url := flag.String("url", "fishioon.com", "The server port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)
	authreq := &pb.AuthReq{
		AuthType:  pb.AuthType_NAME_PWD,
		AuthKey:   "fish",
		AuthValue: "123456",
	}
	authres, err := c.Auth(context.Background(), authreq)
	if err != nil {
		log.Fatalf("auth failed: %v", err)
	}
	sid := authres.SessionId
	//log.Printf("session id:%s\n", sid)
	connreq := &pb.ConnectReq{
		SessionId: sid,
	}
	stream, err := c.Connect(context.Background(), connreq)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}

	waitChan := make(chan bool)

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("stream.Recv", err)
				break
			}
			log.Printf("from_uid:%d content:%s\n", msg.FromUid, msg.Content)
		}
		waitChan <- false
	}()
	joinreq := &pb.JoinGroupReq{
		SessionId: sid,
		Url:       *url,
	}
	joinres, err := c.JoinGroup(context.Background(), joinreq)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	groupid := joinres.GroupId
	log.Printf("room id: %d\n", groupid)

	msg := &pb.Msg{
		FromUid:     0,
		DestId:      groupid,
		Content:     "hello",
		ContentType: pb.ContentType_TEXT,
	}

	for i := 0; i < 10; i++ {
		c.Broadcast(context.Background(), msg)
		time.Sleep(time.Second)
	}
	<-waitChan
}
