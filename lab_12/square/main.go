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

	_, err = nc.Subscribe("pipeline.even", func(m *nats.Msg) {
		var msg Message
		json.Unmarshal(m.Data, &msg)

		msg.Value = msg.Value * msg.Value
		data, _ := json.Marshal(msg)
		
		nc.Publish("pipeline.squared", data)
		log.Printf("квадрат обчислено: %d", msg.Value)
	})

	log.Println("сервіс square запущено...")
	select {}
}