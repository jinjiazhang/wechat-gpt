package main

var config Config

type WechatConfig struct {
	Token     string
	AppID     string
	AppSecret string
}

type OpenAIConfig struct {
	ApiKey string
	Model  string
	Proxy  string
}

type AppConfig struct {
	Port    int
	LogFile string
}

type Config struct {
	App    AppConfig
	Wechat WechatConfig
	OpenAI OpenAIConfig
}
