package main

import (
	"fmt"
	"strings"
)

func UsageText() string {
	return "-reset: reset your conversation\n" +
		"\tc: Chat, f: Funny, k: Keywords, q: Q&A\n" +
		"-model: set openai model\n" +
		"-name: set your name\n" +
		"-friend: set friend's name\n" +
		"-proem: set conversation proem\n"
}

func TextMessage(openid string, text string) (string, error) {
	if text[0] == '-' {
		return CommandMessage(openid, text)
	}

	if openid == "Admin" {
		return AdminMessage(openid, text)
	}

	return ChatMessage(openid, text)
}

func CommandMessage(openid string, text string) (string, error) {
	args := strings.Split(text, " ")
	text = text[len(args)+1:]
	switch args[0] {
	case "-help":
		return HelpCommand(openid, text)
	case "-reset":
		return ResetCommand(openid, text)
	case "-model":
		return ModelCommand(openid, text)
	case "-name":
		return NameCommand(openid, text)
	case "-friend":
		return FriendCommand(openid, text)
	case "-proem":
		return ProemCommand(openid, text)
	default:
		return "", fmt.Errorf("Unknow Command: %s", args[0])
	}
}

func HelpCommand(openid string, text string) (string, error) {
	return UsageText(), nil
}

func ResetCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.Reset(text)
	return "reset mode: " + session.mode, nil
}

func ModelCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.model = text
	return "set model: " + text, nil
}

func NameCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.name = text
	return "set name: " + text, nil
}

func FriendCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.friend = text
	return "set friend: " + text, nil
}

func ProemCommand(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.proem = text + "\n\n"
	return "set proem: " + text, nil
}

func ChatMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.Chat(text)
	return "", nil
}

func AdminMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	return session.Ask(text)
}
