package main

import "net/http"

type Config struct {
	ListenAddr string
	StoreProducerFunc
}

type Server struct {
	*Config
	topics map[string]Storer
}

func NewServer(cfg *Config) (*Server, error) {
	return &Server{
		Config: cfg,
		topics: make(map[string]Storer),
	}, nil
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.ListenAddr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) createTopic(name string) bool {
	if _, ok := s.topics[name]; !ok {
		s.topics[name] = s.StoreProducerFunc()
		return true
	}
	return false
}