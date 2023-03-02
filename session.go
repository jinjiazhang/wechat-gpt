package main

import (
	"sync"
	"time"
)

var (
	sessionMap = make(map[string]*Session)
	mapMutex   sync.Mutex
)

type Session struct {
	openid   string
	mode     string
	model    string
	prologue string
	messages []*ChatGPTMessage
	mutex    sync.Mutex
}

func GetSession(openid string) *Session {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if session, ok := sessionMap[openid]; ok {
		return session
	}

	session := &Session{openid: openid}
	session.Reset("AI")

	sessionMap[openid] = session
	return session
}

func (s *Session) Ask(text string) (string, error) {
	s.push(&ChatGPTMessage{
		Role:    "user",
		Content: text,
		Time:    time.Now().Unix(),
	})

	reply, err := RequestChatGPT(s.model, s.prompt())
	if err != nil {
		return "", err
	}

	s.push(&ChatGPTMessage{
		Role:    "assistant",
		Content: reply,
		Time:    time.Now().Unix(),
	})
	return reply, nil
}

func (s *Session) Chat(text string) error {
	s.push(&ChatGPTMessage{
		Role:    "user",
		Content: text,
		Time:    time.Now().Unix(),
	})

	go s.process(s.prompt())
	return nil
}

func (s *Session) process(messages []*ChatGPTMessage) {
	reply, err := RequestChatGPT(s.model, messages)
	if err != nil {
		SendTextMessage(s.openid, err.Error())
		return
	}

	s.push(&ChatGPTMessage{
		Role:    "assistant",
		Content: reply,
		Time:    time.Now().Unix(),
	})
	SendTextMessage(s.openid, reply)
}

func (s *Session) Reset(mode string) error {
	s.messages = make([]*ChatGPTMessage, 0)
	if mode == "" {
		mode = s.mode
	}

	switch mode {
	case "a", "A", "AI":
		s.mode = "AI"
		s.model = "gpt-3.5-turbo"
		s.prologue = "You are a helpful assistant."
	case "t", "T", "Translate":
		s.mode = "Translate"
		s.model = "gpt-3.5-turbo"
		s.prologue = "You are a helpful assistant that mutual translation between Chinese and English."
	}
	return nil
}

func (s *Session) push(message *ChatGPTMessage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.messages == nil {
		return nil
	}

	s.messages = append(s.messages, message)
	return nil
}

func (s *Session) prompt() []*ChatGPTMessage {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	totalSize := len(s.messages)
	lastTime := time.Now().Unix()
	messages := make([]*ChatGPTMessage, 0)
	tokenCount := 0

	for i := 0; i < totalSize; i++ {
		message := s.messages[totalSize-i-1]
		if lastTime-message.Time > 300 {
			break
		}

		tokenCount += len(message.Content)
		if tokenCount > 2048 {
			break
		}

		lastTime = message.Time
		messages = append([]*ChatGPTMessage{message}, messages...)
	}

	if s.prologue != "" {
		message := &ChatGPTMessage{
			Role:    "system",
			Content: s.prologue,
			Time:    lastTime,
		}
		messages = append([]*ChatGPTMessage{message}, messages...)
	}

	return messages
}
