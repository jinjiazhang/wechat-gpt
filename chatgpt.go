package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ChatGPTRequest struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        uint    `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type ChatGPTResponse struct {
	Error   ErrorItem              `json:"error"`
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ErrorItem struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer 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"
// -d '{"model": "text-davinci-003", "prompt": "地球为什么绕太阳转", "temperature": 0, "max_tokens": 2048}'
func RequestChatGPT(model string, prompt string) (string, error) {
	req := ChatGPTRequest{
		Model:            model,
		Prompt:           prompt,
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

	url := "https://api.openai.com/v1/completions"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)
	client := &http.Client{Timeout: 30 * time.Second}
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
		reply = rsp.Choices[0].Text
	}
	return reply, nil
}
