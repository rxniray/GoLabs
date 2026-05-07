package main 

import (
	"fmt"
	"log"
	"lab_09/serializer"
)

func main() {
	serverConfig := serializer.Server{
		Host:       "127.0.0.1",
		Port:       9090,
		Debug:      false,
		AllowedIPs: []string{"192.168.0.1", "10.0.0.2"},
	}

	fmt.Println("--- Результат у YAML ---")
	yamlData, err := serializer.ToYAML(serverConfig)
	if err != nil {
		log.Fatalf("Помилка генерації YAML: %v", err)
	}
	fmt.Println(yamlData)

	fmt.Println("--- Результат у JSON ---")
	jsonData, err := serializer.ToJSON(serverConfig)
	if err != nil {
		log.Fatalf("Помилка генерації JSON: %v", err)
	}
	fmt.Println(jsonData)
}