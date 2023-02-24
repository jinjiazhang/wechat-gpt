package main

import (
	"fmt"
	"strings"
)

func UsageText() string {
	return "welcome"
}

func TextMessage(openid string, text string) (string, error) {
	if text[0] == '#' {
		return CommandMessage(openid, text)
	}
	return ChatMessage(openid, text)
}

func CommandMessage(openid string, text string) (string, error) {
	args := strings.Split(text, " ")
	switch args[0] {
	case "help":
		return HelpCommand(openid, text)
	case "reset":
		return ResetCommand(openid, text)
	default:
		return "", fmt.Errorf("Unknow Command: %s", args[0])
	}
}

func HelpCommand(openid string, text string) (string, error) {
	return UsageText(), nil
}

func ResetCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.Reset()
	return "reset ok", nil
}

func ChatMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.Chat(text)
	return "", nil
}
