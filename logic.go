package main

import (
	"context"
	"time"
)

func WeChatMessage(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	reply, err := RequestChatGPT(req.Content)
	if err != nil {
		reply = err.Error()
	}

	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      reply,
		Content:      req.Content,
	}

	return rsp, nil
}
