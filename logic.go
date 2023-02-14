package main

import (
	"context"
	"encoding/xml"
	"time"

	log "github.com/sirupsen/logrus"
)

type MessageReq struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        int64  `xml:"MsgId"`
	MsgDataId    int64  `xml:"MsgDataId"`
	Idx          int64  `xml:"Idx"`
	Event        string `xml:"Event"`
}

type MessageRsp struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

func WeChatMessage(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	reply, err := RequestChatGPT(req.Content)
	if err != nil {
		log.Errorf("WeChatMessage request fail, err: %+v", err)
		reply = err.Error()
	}

	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      reply,
	}

	return rsp, nil
}
