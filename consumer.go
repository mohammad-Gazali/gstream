package main

import "github.com/gorilla/websocket"

func Foo() {
	websocket.DefaultDialer.Dial("ws://foo", nil)
}

type Consumer interface {
	Start() error
}