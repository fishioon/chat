package chat

import (
	pb "github.com/fishioon/chat/proto"
)

type ChatGroup struct {
	Id       uint64
	Url      string
	sessions map[string]*Session
	msg      chan *pb.Msg
	join     chan *Session
	leave    chan *Session
}

func NewChatGroup(id uint64, url string) *ChatGroup {
	group := &ChatGroup{
		Id:       id,
		Url:      url,
		msg:      make(chan *pb.Msg),
		join:     make(chan *Session),
		leave:    make(chan *Session),
		sessions: make(map[string]*Session),
	}
	group.listen()
	return group
}

func (group *ChatGroup) listen() {
	go func() {
		for {
			select {
			case msg := <-group.msg:
				for _, session := range group.sessions {
					session.MsgSend(msg)
				}
			case session := <-group.join:
				group.sessions[session.Sid] = session
			case session := <-group.leave:
				delete(group.sessions, session.Sid)
			}
		}
	}()
}

func (group *ChatGroup) RemoveSession(session *Session) {
	group.leave <- session
}

func (group *ChatGroup) AddSession(session *Session) {
	group.join <- session
}

func (group *ChatGroup) RecvBroadcast(msg *pb.Msg) {
	group.msg <- msg
}
