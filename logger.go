package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type LogFormatter struct {
}

func (*LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	if len(entry.Message) == 1 && entry.Message == "\n" {
		return nil, nil
	}

	timeText := entry.Time.Format("15:04:05")
	output := fmt.Sprintf("%s %s\n", timeText, entry.Message)
	return []byte(output), nil
}
