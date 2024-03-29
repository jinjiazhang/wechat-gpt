package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/singleflight"

	log "github.com/sirupsen/logrus"
)

var (
	accessToken = ""
	expiresTime = int64(0)
	singleGroup singleflight.Group
)

type WxTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresTime int64  `json:"expires_in"`
}

type WxNotifyRequest struct {
	ToUser  string   `json:"touser"`
	MsgType string   `json:"msgtype"`
	Text    TextData `json:"text"`
}

type WxNotifyResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func getAccessToken() string {
	if accessToken != "" && time.Now().Unix() < expiresTime {
		return accessToken
	}

	val, _, _ := singleGroup.Do("AccessToken", func() (interface{}, error) {
		return reqAccessToken()
	})

	return val.(string)
}

func reqAccessToken() (string, error) {
	format := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	url := fmt.Sprintf(format, config.Wechat.AppID, config.Wechat.AppSecret)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("GetAccessToken NewRequest fail, err: %+v", err)
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Errorf("GetAccessToken Do fail, err: %+v", err)
		return "", err
	}

	defer response.Body.Close()
	rspBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("GetAccessToken ReadAll fail, err: %+v", err)
		return "", err
	}

	log.Infof("GetAccessToken url: %s, rsp: %s", url, string(rspBody))
	rsp := &WxTokenResponse{}
	err = json.Unmarshal(rspBody, rsp)
	if err != nil {
		log.Errorf("GetAccessToken Unmarshal fail, err: %+v", err)
		return "", err
	}

	if rsp.ErrCode != 0 {
		log.Errorf("GetAccessToken fail, code: %d, msg: %s", rsp.ErrCode, rsp.ErrMsg)
		return "", fmt.Errorf("code: %d, msg: %s", rsp.ErrCode, rsp.ErrMsg)
	}

	accessToken = rsp.AccessToken
	expiresTime = time.Now().Unix() + rsp.ExpiresTime - 600
	return accessToken, nil
}

func sendTextMessage(openid string, text string) error {
	req := &WxNotifyRequest{
		ToUser:  openid,
		MsgType: kMsgTypeText,
		Text: TextData{
			Content: text,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	format := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
	url := fmt.Sprintf(format, getAccessToken())
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Errorf("SendTextMessage NewRequest fail, err: %+v", err)
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Errorf("SendTextMessage Do fail, err: %+v", err)
		return err
	}

	defer response.Body.Close()
	rspBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("SendTextMessage ReadAll fail, err: %+v", err)
		return err
	}

	log.Infof("SendTextMessage req: %s, rsp: %s", string(reqBody), string(rspBody))
	rsp := &WxNotifyResponse{}
	err = json.Unmarshal(rspBody, rsp)
	if err != nil {
		log.Errorf("SendTextMessage Unmarshal fail, err: %+v", err)
		return err
	}

	if rsp.ErrCode != 0 {
		log.Errorf("SendTextMessage fail, code: %d, msg: %s", rsp.ErrCode, rsp.ErrMsg)
		return fmt.Errorf("code: %d, msg: %s", rsp.ErrCode, rsp.ErrMsg)
	}

	return nil
}

func SendTextMessage(openid string, text string) error {
	log.Infof("SendTextMessage openid: %s, text: %s", openid, text)
	segments := splitMessage(text, "\n\n", 2048)
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		err := sendTextMessage(openid, segment)
		if err != nil {
			sendTextMessage(openid, err.Error())
		}
	}

	return nil
}

func splitMessage(text string, sep string, maxLen int) []string {
	if len(text) < maxLen {
		return []string{text}
	}

	parts := strings.Split(text, sep)
	result := make([]string, 0)
	var temp string
	for _, part := range parts {
		if len(temp)+len(part) < maxLen {
			if temp != "" {
				temp += sep
			}
			temp += part
		} else {
			result = append(result, temp)
			if len(part) > maxLen {
				size := len(part) / maxLen
				for i := 0; i < size; i++ {
					result = append(result, part[i*maxLen:(i+1)*maxLen])
				}
				part = part[size*maxLen:]
			}

			temp = part
		}
	}

	if temp != "" {
		result = append(result, temp)
	}
	return result
}
