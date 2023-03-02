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

func CommandMessage(openid string, text string) (string, error) {
	args := strings.Split(text, " ")
	if len(text) > len(args[0])+1 {
		text = text[len(args[0])+1:]
	}

	switch args[0] {
	case "-help":
		return HelpCommand(openid, text)
	case "-reset":
		return ResetCommand(openid, text)
	case "-model":
		return ModelCommand(openid, text)
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
