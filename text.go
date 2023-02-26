package main

func TextMessage(openid string, text string) (string, error) {
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
	session.Chat(text)
	return "", nil
}

func AdminMessage(openid string, text string) (string, error) {
	session := GetSession(openid)
	return session.Ask(text)
}

func EventMessage(openid string, event string, eventKey string) (string, error) {
	switch event {
	case "CLICK":
		return MenuCommand(openid, eventKey)
	default:
		return "I'm the AI assistant make by Jinjiazh, Let's start our conversation!", nil
	}
}

func MenuCommand(openid string, menuKey string) (string, error) {
	session := GetSession(openid)
	switch menuKey {
	case "TEXT_AI":
		session.Reset("AI")
	case "TEXT_CHAT":
		session.Reset("CHAT")
	case "TEXT_QA":
		session.Reset("QA")
	case "TEXT_FUNNY":
		session.Reset("FUNNY")
	case "TEXT_RESET":
		session.Reset("")
	case "IMAGE_DAllE":
		return "功能暂未开放", nil
	case "BUTTON_ABOUT":
		return "此公众号基于个人兴趣开发，使用OpenAI接口生成回应内容。仅供学习用途，严禁他用！", nil
	default:
		return "功能暂未开放", nil
	}

	return session.proem, nil
}
