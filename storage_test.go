package main

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {

	s := NewMemoryStorage()
	
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("foobaz-%d", i)

		latestOffset, err := s.Push([]byte(key))
		
		if err != nil {
			t.Error(err)
		}

		data, err := s.Fetch(latestOffset)

		if err != nil {
			t.Error(err)
		}

		fmt.Println(string(data))
	}
}