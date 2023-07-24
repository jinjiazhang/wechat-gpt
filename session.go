package main

import (
	"fmt"
	"regexp"
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
	messages []*ChatGPTMessage
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
		mode:     "AI",
		model:    config.OpenAI.Model,
		messages: make([]*ChatGPTMessage, 0),
	}

	sessionMap[openid] = session
	return session
}

func (s *Session) SyncAsk(text string) (string, error) {
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

func (s *Session) AsyncAsk(text string) error {
	s.push(&ChatGPTMessage{
		Role:    "user",
		Content: text,
		Time:    time.Now().Unix(),
	})

	go func() {
		reply, err := RequestChatGPT(s.model, s.prompt())
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
	}()
	return nil
}

func (s *Session) Translate(text string) error {
	result, _ := regexp.MatchString(`[\x{4e00}-\x{9fa5}]+`, text)
	format := "Translate the following English text to Chinese:\"%s\""
	if result {
		format = "Translate the following Chinese text to English:\"%s\""
	}

	message := &ChatGPTMessage{
		Role:    "user",
		Content: fmt.Sprintf(format, text),
		Time:    time.Now().Unix(),
	}

	go func() {
		reply, err := RequestChatGPT(s.model, []*ChatGPTMessage{message})
		if err != nil {
			SendTextMessage(s.openid, err.Error())
			return
		}

		SendTextMessage(s.openid, reply)
	}()
	return nil
}

func (s *Session) Process(text string) error {
	switch s.mode {
	case "t", "Translate":
		return s.Translate(text)
	default:
		return s.AsyncAsk(text)
	}
}

func (s *Session) Reset(mode string) error {
	s.messages = make([]*ChatGPTMessage, 0)
	if mode == "" {
		mode = s.mode
	}

	switch mode {
	case "a", "AI":
		s.mode = "AI"
	case "t", "Translate":
		s.mode = "Translate"
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
		if lastTime-message.Time > 600 {
			break
		}

		tokenCount += int(message.TokenNum)
		if tokenCount > 2048 {
			break
		}

		lastTime = message.Time
		messages = append([]*ChatGPTMessage{message}, messages...)
	}

	message := &ChatGPTMessage{
		Role:     "system",
		Content:  "You are a helpful assistant.",
		TokenNum: 13,
		Time:     lastTime,
	}
	messages = append([]*ChatGPTMessage{message}, messages...)
	return messages
}
