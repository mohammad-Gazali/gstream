package main

import (
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Peer interface {
	Send([]byte) error
}

type WSPeer struct {
	conn   *websocket.Conn
	server *Server
}

func NewWSPeer(conn *websocket.Conn, server *Server) *WSPeer {
	p := &WSPeer{
		conn:   conn,
		server: server,
	}

	go p.readLoop()

	return p
}

func (p *WSPeer) readLoop() {
	var msg WSMessage
	for {
		if err := p.conn.ReadJSON(&msg); err != nil {
			slog.Error("ws peer read error", "err", err)
			return
		}

		if err := p.handleMessage(msg); err != nil {
			slog.Error("ws peer handle message error", "err", err)
			return
		}

	}
}

func (p *WSPeer) Send(data []byte) error {
	return p.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (p *WSPeer) handleMessage(msg WSMessage) error {
	if msg.Action == "subscribe" {
		p.server.AddPeerToTopics(p, msg.Topics)
	}
	fmt.Printf("handling message => %+v\n", msg)
	return nil
}
