package main

import "encoding/xml"

var (
	kMsgTypeEvent = "event"
	kMsgTypeText  = "text"
	kMsgTypeImage = "image"
)

type EventData struct {
	Event    string `xml:"Event"`
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`
}

type TextData struct {
	Content string `xml:"Content"`
}

type ImageData struct {
	PicUrl  string `xml:"PicUrl"`
	MediaId string `xml:"MediaId"`
}

type MessageReq struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	EventData
	TextData
	ImageData
	MsgId     int64 `xml:"MsgId"`
	MsgDataId int64 `xml:"MsgDataId"`
	Idx       int64 `xml:"Idx"`
}

type MessageRsp struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	TextData
	ImageData
}
