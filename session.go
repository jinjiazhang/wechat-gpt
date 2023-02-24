package main

import (
	"sync"
	"time"
)

var (
	sessionMap map[string]*Session
	mapMutex   sync.Mutex
)

type Message struct {
	author string
	text   string
	time   int64
}

type Session struct {
	openid   string
	messages []*Message
	mutex    sync.Mutex
}

func GetSession(openid string) *Session {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if session, ok := sessionMap[openid]; ok {
		return session
	}

	session := &Session{
		openid:   openid,
		messages: make([]*Message, 0),
	}
	sessionMap[openid] = session
	return session
}

func (s *Session) Chat(text string) error {
	s.push(&Message{
		author: "Human",
		text:   text,
		time:   time.Now().Unix(),
	})

	go func() {
		reply, err := RequestChatGPT(s.prompt())
		if err != nil {
			NotifyUser(s.openid, err.Error())
			return
		}

		s.push(&Message{
			author: "Robot",
			text:   reply,
			time:   time.Now().Unix(),
		})
		NotifyUser(s.openid, reply)
	}()

	return nil
}

func (s *Session) Reset() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = make([]*Message, 0)
	return nil
}

func (s *Session) push(message *Message) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = append(s.messages, message)
	return nil
}

func (s *Session) prompt() string {
	return ""
}
