package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type Consumer interface {
	Start() error
}

var upgrader = websocket.Upgrader{}

type WSConsumer struct {
	listenAddr string
	server     *Server
}

func NewWSConsumer(listenAddr string, server *Server) *WSConsumer {
	return &WSConsumer{
		listenAddr: listenAddr,
		server:     server,
	}
}

func (ws *WSConsumer) Start() error {
	slog.Info("Websocket consumer started", "port", ws.listenAddr)
	return http.ListenAndServe(ws.listenAddr, ws)
}

func (ws *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := NewWSPeer(conn, ws.server)
	ws.server.AddPeer(p)
}

type WSMessage struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}