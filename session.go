package main

type Session struct {
}

func (s *Session) Chat(text string) error {
	return nil
}

func (s *Session) Reset() error {
	return nil
}

func GetSession(openid string) *Session {
	return nil
}
