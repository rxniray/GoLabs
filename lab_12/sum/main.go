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

	totalSum := 0

	_, err = nc.Subscribe("pipeline.squared", func(m *nats.Msg) {
		var msg Message
		json.Unmarshal(m.Data, &msg)

		totalSum += msg.Value
		log.Printf("отримано: %d | поточна підсумкова сума: %d", msg.Value, totalSum)
	})

	log.Println("сервіс sum запущено...")
	select {}
}