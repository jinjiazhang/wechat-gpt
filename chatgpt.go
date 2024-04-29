package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type ChatGPTMessage struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	TokenNum int32  `json:"-"`
	Time     int64  `json:"-"`
}

type ChatGPTRequest struct {
	Model            string            `json:"model"`
	Messages         []*ChatGPTMessage `json:"messages"`
	MaxTokens        uint              `json:"max_tokens"`
	Temperature      float64           `json:"temperature"`
	TopP             int               `json:"top_p"`
	FrequencyPenalty int               `json:"frequency_penalty"`
	PresencePenalty  int               `json:"presence_penalty"`
}

type ChatGPTResponse struct {
	Error   ChatGPTError    `json:"error"`
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int             `json:"created"`
	Usage   ChatGPTUsage    `json:"usage"`
	Choices []ChatGPTChoice `json:"choices"`
}

type ChatGPTError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ChatGPTChoice struct {
	Index        int             `json:"index"`
	FinishReason string          `json:"finish_reason"`
	Message      *ChatGPTMessage `json:"message"`
}

type ChatGPTUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// curl https://api.openai.com/v1/chat/completions \
//   -H 'Content-Type: application/json' \
//   -H 'Authorization: Bearer YOUR_API_KEY' \
//   -d '{
//   "model": "gpt-3.5-turbo",
//   "messages": [{"role": "user", "content": "Hello!"}]
// }'
func RequestChatGPT(model string, messages []*ChatGPTMessage) (string, error) {
	req := ChatGPTRequest{
		Model:            model,
		Messages:         messages,
		MaxTokens:        2048,
		Temperature:      0.8,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	chatUrl := "https://api.openai.com/v1/chat/completions"
	request, err := http.NewRequest("POST", chatUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.OpenAI.ApiKey)
	client := &http.Client{Timeout: 300 * time.Second}

	if config.OpenAI.Proxy != "" {
		proxyUrl, _ := url.Parse(config.OpenAI.Proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("http error, code: %d, detail: %s", response.StatusCode, string(body))
	}

	rspBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	log.Infof("RequestChatGPT req: %s, rsp: %s", string(reqBody), string(rspBody))
	rsp := &ChatGPTResponse{}
	err = json.Unmarshal(rspBody, rsp)
	if err != nil {
		return "", err
	}

	if rsp.Error.Message != "" {
		return "", fmt.Errorf("chatgpt error, type: %s, message: %s", rsp.Error.Type, rsp.Error.Message)
	}

	var reply string
	if len(rsp.Choices) > 0 {
		reply = rsp.Choices[0].Message.Content
	}
	return reply, nil
}

func ProxyChatGPT(w http.ResponseWriter, r *http.Request) {
	chatUrl := "https://api.openai.com/v1/chat/completions"
	request, err := http.NewRequest(r.Method, chatUrl, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Header.Set("Authorization", "Bearer "+config.OpenAI.ApiKey)
	for header, values := range r.Header {
		for _, value := range values {
			request.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制目标服务器的响应头部和响应体到原始响应中
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
