package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	_, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	playerName := "Player"
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}

	log.Printf("Player: %s has joined to game", playerName)
}

func main() {
	log.Println("wssrv starting...")

	http.HandleFunc("/", wsHandler)

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
