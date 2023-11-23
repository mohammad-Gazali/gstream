package main

import (
	"fmt"
	"log/slog"
	"sync"
)

type Message struct {
	Topic string
	Data  []byte
}

type Config struct {
	HTTPListenAddr string
	WSListenAddr   string
	StoreProducerFunc
}

type Server struct {
	*Config

	topics map[string]Storer

	mu    sync.RWMutex
	peers map[Peer]bool

	consumers []Consumer
	producers []Producer

	producech chan Message
	quitch    chan struct{}
}

func NewServer(cfg *Config) (*Server, error) {
	producech := make(chan Message)
	s := &Server{
		Config:    cfg,
		topics:    make(map[string]Storer),
		quitch:    make(chan struct{}),
		producech: producech,
		producers: []Producer{
			NewHTTPProducer(cfg.HTTPListenAddr, producech),
		},
		mu:    sync.RWMutex{},
		peers: make(map[Peer]bool),
		consumers: []Consumer{},
	}

	s.consumers = append(s.consumers, NewWSConsumer(cfg.WSListenAddr, s))

	return s, nil
}

func (s *Server) Start() {
	for _, c := range s.consumers {
		go func(c Consumer) {
			if err := c.Start(); err != nil {
				fmt.Println(err)
			}
		}(c)
	}

	for _, p := range s.producers {
		go func(p Producer) {
			if err := p.Start(); err != nil {
				fmt.Println(err)
			}
		}(p)
	}

	s.loop()
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitch:
			return
		case msg := <-s.producech:
			offset, err := s.publish(msg)
			if err != nil {
				slog.Error("failed to publish", "error", err)
			} else {
				slog.Info("produced message", "offset", offset)
			}
		}
	}
}

func (s *Server) publish(msg Message) (int, error) {
	store := s.getStoreForTopic(msg.Topic)
	return store.Push(msg.Data)
}

func (s *Server) getStoreForTopic(topic string) Storer {
	if _, ok := s.topics[topic]; !ok {
		s.topics[topic] = s.StoreProducerFunc()
		slog.Info("created new topic", "topic", topic)
	}
	return s.topics[topic]
}

func (s *Server) AddPeer(p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("added new peer", "peer", p)
	s.peers[p] = true
}

func (s *Server) AddPeerToTopics(p Peer, topics []string) {
	// send all the messages from the peer's offset
	for _, topic := range topics {
		store := s.getStoreForTopic(topic)
		// size := store.Len()
		// for i := 0; i < size; i++ {
		// 	b, _ := store.Get(i)
			
		// }
	}
	slog.Info("adding peer to topics", "topics", topics, "peer", p)
}