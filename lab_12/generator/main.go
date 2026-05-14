package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Value int `json:"value"`
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("помилка підключення до NATS: %v", err)
	}
	defer nc.Close()

	time.Sleep(3 * time.Second)
	log.Println("генератор починає роботу...")

	for i := 1; i <= 100; i++ {
		msg := Message{Value: i}
		data, _ := json.Marshal(msg)

		err := nc.Publish("pipeline.numbers", data)
		if err != nil {
			log.Printf("помилка публікації: %v", err)
		}
	}
	
	nc.Flush() 
	log.Println("генератор завершив публікацію.")
}