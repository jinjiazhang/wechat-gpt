package main

import log "github.com/sirupsen/logrus"

func TextMessage(openid string, text string) (string, error) {
	log.Infof("RecvTextMessage openid: %s, text: %s", openid, text)
	if text[0] == '-' {
		return CommandMessage(openid, text)
	}

	if openid == "Admin" {
		return AdminMessage(openid, text)
	}

	return ChatMessage(openid, text)
}

func ChatMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	session.Process(text)
	return "", nil
}

func AdminMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	return session.SyncAsk(text)
}

func EventMessage(openid string, event string, eventKey string) (string, error) {
	switch event {
	case "CLICK":
		return MenuCommand(openid, eventKey)
	default:
		return "有什么可以帮助你吗", nil
	}
}

func MenuCommand(openid string, menuKey string) (string, error) {
	session := GetSession(openid)
	switch menuKey {
	case "TEXT_AI":
		session.Reset("AI")
		return "有什么可以帮助你吗", nil
	case "TEXT_TRANSLATE":
		session.Reset("Translate")
		return "请输入需要翻译的内容", nil
	case "TEXT_RESET":
		session.Reset("")
		return "对话已重置", nil
	case "IMAGE_DAllE":
		return "功能暂未开放", nil
	case "BUTTON_ABOUT":
		return "此公众号基于个人兴趣开发，使用OpenAI接口生成回应内容。仅供学习用途，严禁他用！", nil
	default:
		return "功能暂未开放", nil
	}
}
