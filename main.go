package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"

	log "github.com/sirupsen/logrus"
)

func main() {
	apiKey := flag.String("key", "sk-***", "chatgpt api-key")
	token := flag.String("token", "jinjiazh", "wechat message token")
	appId := flag.String("appid", "wxf963***", "wechat appid")
	appSecret := flag.String("secret", "bf8fd***", "wechat app secret")
	flag.Parse()

	OPENAI_API_KEY = *apiKey
	WECHAT_TOKEN = *token
	WECHAT_APPID = *appId
	WECHAT_APPSECRET = *appSecret
	http.HandleFunc("/message", HandleMessage)
	http.ListenAndServe(":8080", nil)
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		HandleMessage_GET(w, r)
	case "POST":
		HandleMessage_POST(w, r)
	case "PUT":
		HandleMessage_PUT(w, r)
	}
}

func HandleMessage_GET(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url: %s", r.URL.String())
	query := r.URL.Query()

	signature := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")

	sl := []string{WECHAT_TOKEN, timestamp, nonce}
	sort.Strings(sl)
	sum := sha1.Sum([]byte(sl[0] + sl[1] + sl[2]))
	if signature == hex.EncodeToString(sum[:]) {
		w.Write([]byte(echostr))
	}
}

func HandleMessage_POST(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("HandleMessage read body fail, err: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	req := &MessageReq{}
	err = xml.Unmarshal(reqBody, req)
	if err != nil {
		log.Errorf("HandleMessage unmarshal req fail, err: %+v", err)
		return
	}

	rsp, err := WeChatMessage(context.TODO(), req)
	if err != nil {
		log.Errorf("HandleMessage handle func fail, err: %+v", err)
		return
	}

	if rsp == nil {
		// 回复success，这样微信服务器不会对此作任何处理，并且不会发起重试
		w.Write([]byte("success"))
		return
	}

	rspBody, err := xml.Marshal(rsp)
	if err != nil {
		log.Errorf("HandleMessage marshal rsp fail, err: %+v", err)
		return
	}

	_, err = w.Write(rspBody)
	if err != nil {
		log.Errorf("HandleMessage write body fail, err: %+v", err)
		return
	}

	log.Infof("WeChatMessage req: %s, rsp: %s", string(reqBody), string(rspBody))
}

func HandleMessage_PUT(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("HandleMessage read body fail, err: %+v", err)
		return
	}

	reply, err := TextMessage("Admin", string(reqBody))
	if err != nil {
		log.Errorf("HandleMessage err: %+v", err)
		reply = err.Error()
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(reply))
	if err != nil {
		log.Errorf("HandleMessage write body fail, err: %+v", err)
		return
	}

	log.Infof("HandleMessage req: %s, rsp: %s", string(reqBody), reply)
}
