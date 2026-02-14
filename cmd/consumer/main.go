package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/example/ms-golang-event-driven/internal/rabbitmq"
)

const queueName = "training.events"

type trainingMessage struct {
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	client := rabbitmq.NewClient()
	if err := client.DeclareQueue(queueName); err != nil {
		log.Fatalf("falha ao declarar fila: %v", err)
	}

	log.Println("consumer aguardando mensagens...")
	for {
		body, err := client.GetOne(queueName)
		if err != nil {
			log.Printf("erro ao buscar mensagem: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if body == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		var payload trainingMessage
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("mensagem invalida: %s", string(body))
			continue
		}

		log.Printf("[EVENTO] name=%s message=%s created_at=%s", payload.Name, payload.Message, payload.CreatedAt.Format(time.RFC3339))
	}
}
