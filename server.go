package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slack-bot/middleware"
)

var addr = ":8080"

type event struct {
	Channel     string `json:"channel"`
	ChannelType string `json:"channel_type"`
	ClientMsgId string `json:"client_msg_id"`
	EventTs     int    `json:"event_ts"` // unix time
	Team        string `json:"team"`
	Text        string `json:"text"`
	Ts          int    `json:"ts"` // unix time
	Type        string `json:"type"`
	User        string `json:"user"` // 例のID
}

type req struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Event     event  `json:"event"`
}

func connectionChallenge(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	var request req
	json.Unmarshal(buf.Bytes(), &request)

	log.Println(request)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, request.Challenge)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/connection-challenge", connectionChallenge)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Health Check\n")
	})
	fmt.Printf("[START] server. port: %s\n", addr)

	if err := http.ListenAndServe(addr, middleware.Log(router)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
