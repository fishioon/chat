package main

import (
	"errors"
	"github.com/fishioon/chat"
	pb "github.com/fishioon/chat/proto"
	"golang.org/x/net/context"
	"hash/fnv"
	"log"
)

type chatServer struct {
	groups   map[uint64]*chat.ChatGroup
	sessions map[string]*chat.Session
	users    *chat.UserManager
}

func hashUrl(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func NewChatServer() *chatServer {
	return &chatServer{
		groups:   make(map[uint64]*chat.ChatGroup),
		sessions: make(map[string]*chat.Session),
		users:    chat.NewUserManager(),
	}
}

func (cs *chatServer) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthRes, error) {
	switch req.AuthType {
	case pb.AuthType_EMAIL_PWD:
	case pb.AuthType_PHONE_PWD:
	case pb.AuthType_PHONE_CODE:
	}

	user := cs.users.AuthOrNew(req.AuthKey, req.AuthValue, "")

	session := chat.NewSession(user.Uid)
	cs.sessions[session.Sid] = session

	res := &pb.AuthRes{
		SessionId: session.Sid,
		Uid:       user.Uid,
	}
	log.Println("auth")
	return res, nil
}

func (cs *chatServer) Connect(req *pb.ConnectReq, stream pb.Chat_ConnectServer) error {
	// check session
	session, ok := cs.sessions[req.SessionId]
	if !ok {
		err := errors.New("invalid session id")
		return err
	}
	defer func() {
		session.Offline()
	}()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case msg := <-session.MsgRecv():
			log.Printf("msg send, sid:%s content:%s\n", session.Sid, msg.Content)
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
	log.Println("connect return")
	return nil
}

func (cs *chatServer) Broadcast(ctx context.Context, req *pb.Msg) (*pb.Msg, error) {
	group, ok := cs.groups[req.DestId]
	if !ok {
		log.Printf("invalid group id:%d\b", req.DestId)
		return nil, errors.New("invalid group id")
	}
	log.Printf("broadcast, groupid:%d content:%s\n", req.DestId, req.Content)
	group.RecvBroadcast(req)
	return req, nil
}

func (cs *chatServer) Unicast(ctx context.Context, req *pb.Msg) (*pb.Msg, error) {
	return nil, nil
}

func (cs *chatServer) JoinGroup(ctx context.Context, req *pb.JoinGroupReq) (*pb.JoinGroupRes, error) {
	session, ok := cs.sessions[req.SessionId]
	if !ok {
		log.Printf("invalid session id:%d\b", req.SessionId)
		err := errors.New("invalid session id")
		return nil, err
	}
	log.Printf("join group, url:%s\n", req.Url)
	groupid := hashUrl(req.Url)
	group, ok := cs.groups[groupid]
	if !ok {
		group = chat.NewChatGroup(groupid, req.Url)
		cs.groups[group.Id] = group
	}
	session.JoinGroup(group)
	res := &pb.JoinGroupRes{
		GroupId: group.Id,
	}
	return res, nil
}

func (cs *chatServer) LeaveGroup(ctx context.Context, req *pb.LeaveGroupReq) (*pb.None, error) {
	session, ok := cs.sessions[req.SessionId]
	if !ok {
		err := errors.New("invalid session id")
		return nil, err
	}
	group, ok := cs.groups[req.GroupId]
	if !ok {
		err := errors.New("invalid group id")
		return nil, err
	}
	session.LeaveGroup(group)
	return &pb.None{}, nil
}
