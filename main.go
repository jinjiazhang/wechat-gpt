package main

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/message", handleMessage)
	http.ListenAndServe(":8080", nil)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
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

	rsp, err := HandleMessage(context.TODO(), req)
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

	log.Infof("handleMessage req: %s, rsp: %s", string(reqBody), string(rspBody))
}
