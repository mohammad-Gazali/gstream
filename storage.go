package main

import (
	"fmt"
	"sync"
)

type Storer interface {
	Push([]byte) (int, error)
	Get(int) ([]byte, error)
	// Len() int
}

type StoreProducerFunc func() Storer

type MemoryStorage struct {
	mu   sync.RWMutex
	data [][]byte
}


func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make([][]byte, 0),
	}
}

func (s *MemoryStorage) Push(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

func (s *MemoryStorage) Get(offset int) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if offset < 0 {
		return nil, fmt.Errorf("offset can't be less than zero")
	}

	if len(s.data) <= offset {
		return nil, fmt.Errorf("offset (%d) too high", offset)
	}

	return s.data[offset], nil
}

// func (s *MemoryStorage) Len() int {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	return len(s.data)
// }