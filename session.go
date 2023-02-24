package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	sessionMap = make(map[string]*Session)
	mapMutex   sync.Mutex
)

type Message struct {
	author string
	text   string
	time   int64
}

type Session struct {
	openid   string
	name     string
	proem    string
	friend   string
	messages []*Message
	mutex    sync.Mutex
}

func GetSession(openid string) *Session {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if session, ok := sessionMap[openid]; ok {
		return session
	}

	session := &Session{openid: openid}
	session.Reset()

	sessionMap[openid] = session
	return session
}

func (s *Session) Ask(text string) (string, error) {
	s.push(&Message{
		author: s.name,
		text:   text,
		time:   time.Now().Unix(),
	})

	reply, err := RequestChatGPT(s.prompt())
	if err != nil {
		return "", err
	}

	s.push(&Message{
		author: s.friend,
		text:   reply,
		time:   time.Now().Unix(),
	})
	return reply, nil
}

func (s *Session) Chat(text string) error {
	s.push(&Message{
		author: s.name,
		text:   text,
		time:   time.Now().Unix(),
	})

	go s.process(s.prompt())
	return nil
}

func (s *Session) Reset() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.name = "Human"
	s.friend = "AI"
	s.proem = "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n\n"
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	text := fmt.Sprintf("%s:", s.friend)
	size := len(s.messages)
	for i := 0; i < size; i++ {
		message := s.messages[size-i-1]
		if time.Now().Unix()-message.time > 7200 {
			break
		}

		text = fmt.Sprintf("%s: %s\n%s", message.author, message.text, text)
	}
	return s.proem + text
}

func (s *Session) process(prompt string) {
	reply, err := RequestChatGPT(prompt)
	if err != nil {
		SendTextMessage(s.openid, err.Error())
		return
	}

	s.push(&Message{
		author: s.friend,
		text:   reply,
		time:   time.Now().Unix(),
	})
	SendTextMessage(s.openid, reply)
}
