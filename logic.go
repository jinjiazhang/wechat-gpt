package main

import (
	"context"
	"time"
)

func HandleMessage(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      req.Content,
	}

	return rsp, nil
}
