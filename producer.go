package main

import (
	"log/slog"
	"net/http"
	"strings"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr string
	producech  chan<- Message
}

func NewHTTPProducer(listenAddr string, producech chan<- Message) *HTTPProducer {
	return &HTTPProducer{
		listenAddr: listenAddr,
		producech: producech,
	}
}

func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path  = strings.TrimPrefix(r.URL.Path, "/")
		parts = strings.Split(path, "/")
	)

	// publishing case
	if r.Method == "POST" {
		if len(parts) != 2 {
			slog.Error("invalid action", "path", r.URL.Path)
			return
		}
		topic := parts[1]
		
		p.producech <- Message{
			Topic: topic,
			Data: []byte("I don't care"),
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	// commit case
	// if r.Method == "GET" {
	// }
}


func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started", "port", p.listenAddr)
	return http.ListenAndServe(p.listenAddr, p)
}
