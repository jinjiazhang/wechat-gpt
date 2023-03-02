package main

import (
	"fmt"
	"strings"
)

func UsageText() string {
	return "usage: [-help] [-reset <mode>]\n" +
		"-reset\treset conversation history"
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
