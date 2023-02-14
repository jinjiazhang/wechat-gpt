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
	apiKey := flag.String("key", "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", "chatgpt api-key")
	flag.Parse()

	API_KEY = *apiKey
	http.HandleFunc("/message", HandleMessage)
	http.ListenAndServe(":8080", nil)
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		HandleMessage_GET(w, r)
	case "POST":
		HandleMessage_POST(w, r)
	}
}

func HandleMessage_GET(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url: %s", r.URL.String())
	query := r.URL.Query()

	token := "jinjiazh"
	signature := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")

	sl := []string{token, timestamp, nonce}
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

	rspBody, err := xml.Marshal(rsp)
	if err != nil {
		log.Errorf("HandleMessage marshal rsp fail, err: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(rspBody)
	if err != nil {
		log.Errorf("HandleMessage write body fail, err: %+v", err)
		return
	}

	log.Infof("WeChatMessage req: %s, rsp: %s", string(reqBody), string(rspBody))
}
