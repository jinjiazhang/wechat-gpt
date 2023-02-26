package main

import (
	"fmt"
	"strings"
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
	model    string
	mode     string
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
	session.Reset("")

	sessionMap[openid] = session
	return session
}

func (s *Session) Ask(text string) (string, error) {
	s.push(&Message{
		author: s.name,
		text:   text,
		time:   time.Now().Unix(),
	})

	reply, err := RequestChatGPT(s.model, s.prompt())
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

func (s *Session) Reset(mode string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.messages = make([]*Message, 0)
	switch strings.ToUpper(mode) {
	case "A", "AI":
		s.mode = "AI assistant"
		s.model = "text-davinci-003"
		s.name = "Human"
		s.friend = "AI"
		s.proem = "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n\n"
	case "C", "CHAT":
		s.mode = "Chat"
		s.model = "text-davinci-003"
		s.name = "You"
		s.friend = "Friend"
		s.proem = ""
	case "Q", "QA":
		s.mode = "Q&A"
		s.model = "text-davinci-003"
		s.name = "Q"
		s.friend = "A"
		s.proem = "I am a highly intelligent question answering bot. If you ask me a question that is rooted in truth, I will give you the answer. If you ask me a question that is nonsense, trickery, or has no clear answer, I will respond with \"Unknown\".\n\n"
	case "F", "FUNNY":
		s.mode = "FUNNY"
		s.model = "text-davinci-003"
		s.name = "You"
		s.friend = "Marv"
		s.proem = "Marv is a chatbot that reluctantly answers questions with humorous responses:\n\n"
	}

	return nil
}

func (s *Session) push(message *Message) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.messages == nil {
		return nil
	}

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

		if len(s.proem)+len(text)+len(message.author)+len(message.text) > 2000 {
			break
		}

		text = fmt.Sprintf("%s: %s\n%s", message.author, message.text, text)
	}
	return s.proem + text
}

func (s *Session) process(prompt string) {
	reply, err := RequestChatGPT(s.model, prompt)
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
