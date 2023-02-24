package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func WeChatMessage(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	switch req.MsgType {
	case kMsgTypeEvent:
		return WeChatEvent(ctx, req)
	case kMsgTypeText:
		return WeChatText(ctx, req)
	case kMsgTypeImage:
		return WeChatImage(ctx, req)
	default:
		rsp := &MessageRsp{
			ToUserName:   req.FromUserName,
			FromUserName: req.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      kMsgTypeText,
			TextData: TextData{
				Content: fmt.Sprintf("Unknow MsgType: %s", req.MsgType),
			},
		}
		return rsp, nil
	}
}

func WeChatEvent(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      kMsgTypeText,
		TextData: TextData{
			Content: "I'm the AI assistant make by Jinjiazh, Let's start our conversation!",
		},
	}
	return rsp, nil
}

func WeChatText(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	reply, err := TextMessage(req.FromUserName, req.Content)
	if err != nil {
		log.Errorf("TextMessage err: %+v", err)
		reply = fmt.Sprintf("TextMessage err: %+v", err)
	}

	if reply == "" {
		return nil, nil
	}

	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      kMsgTypeText,
		TextData: TextData{
			Content: reply,
		},
	}
	return rsp, nil
}

func WeChatImage(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      kMsgTypeImage,
		ImageData: ImageData{
			MediaId: req.ImageData.MediaId,
		},
	}
	return rsp, nil
}
