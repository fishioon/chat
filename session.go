package chat

import (
	pb "github.com/fishioon/chat/proto"
	"log"
	//"math/rand"
)

const (
	msgbufSzie = 10
)

type Session struct {
	Sid    string //session id
	Uid    uint64
	msgbuf chan *pb.Msg
	online bool
	groups map[uint64]*ChatGroup
}

func NewSession(uid uint64) *Session {
	//sid := newSessionID()
	sid := RandStringBytes(32)
	return &Session{
		Sid:    sid,
		Uid:    uid,
		msgbuf: make(chan *pb.Msg, msgbufSzie),
		groups: make(map[uint64]*ChatGroup),
		online: true,
	}
}

func (session *Session) MsgSend(msg *pb.Msg) {
	if session.online {
		session.msgbuf <- msg
	}
}

func (session *Session) MsgRecv() <-chan *pb.Msg {
	return session.msgbuf
}

func (session *Session) Offline() {
	session.online = false
}

func (session *Session) JoinGroup(group *ChatGroup) {
	_, ok := session.groups[group.Id]
	if !ok {
		group.AddSession(session)
		session.groups[group.Id] = group
	} else {
		log.Printf("already join the group:%d\n", group.Id);
		session.online = true;
	}
}

func (session *Session) LeaveGroup(group *ChatGroup) {
	_, ok := session.groups[group.Id]
	if ok {
		group.RemoveSession(session)
		delete(session.groups, group.Id)
	}
}
