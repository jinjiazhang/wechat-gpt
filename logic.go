package main

import (
	"context"
	"fmt"
	"time"
)

func WelcomeText() string {
	return "welcome"
}

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
				Content: fmt.Sprintf("MsgType: %s", req.MsgType),
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
			Content: WelcomeText(),
		},
	}
	return rsp, nil
}

func WeChatText(ctx context.Context, req *MessageReq) (*MessageRsp, error) {
	rsp := &MessageRsp{
		ToUserName:   req.FromUserName,
		FromUserName: req.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      kMsgTypeText,
		TextData: TextData{
			Content: req.Content,
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
