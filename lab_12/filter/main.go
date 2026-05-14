package main

import (
	"encoding/json"
	"log"
	"os"

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
		log.Fatalf("помилка: %v", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("pipeline.numbers", func(m *nats.Msg) {
		var msg Message
		json.Unmarshal(m.Data, &msg)

		if msg.Value%2 == 0 {
			data, _ := json.Marshal(msg)
			nc.Publish("pipeline.even", data)
			log.Printf("відфільтровано: %d", msg.Value)
		}
	})

	log.Println("сервіс filter запущено...")
	select {} 
}