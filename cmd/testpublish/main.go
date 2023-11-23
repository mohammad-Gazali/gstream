package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	url := "http://localhost:3000/publish"
	topics := []string{"topic-1", "topic-2", "topic-3"}
	
	for i := 0; i < 100; i++ {
		topic := topics[rand.Intn(len(topics))]
		payload := []byte(fmt.Sprintf("foobaz-%d", i))
		res, err := http.Post(url + "/" + topic, "application/octet-stream", bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != http.StatusCreated {
			log.Fatal("status code is not 201")
		}
		fmt.Println(res)
	}
}