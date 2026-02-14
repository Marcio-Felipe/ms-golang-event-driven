package main

import (
	"encoding/json"
	"log"
	"net/http"
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

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "metodo nao permitido", http.StatusMethodNotAllowed)
			return
		}

		var payload trainingMessage
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "json invalido", http.StatusBadRequest)
			return
		}

		if payload.Name == "" || payload.Message == "" {
			http.Error(w, "name e message sao obrigatorios", http.StatusBadRequest)
			return
		}

		payload.CreatedAt = time.Now().UTC()
		body, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "erro ao serializar payload", http.StatusInternalServerError)
			return
		}

		if err := client.Publish(queueName, body); err != nil {
			http.Error(w, "erro ao publicar mensagem", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "mensagem publicada"})
	})

	log.Println("producer rodando em :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("erro no servidor http: %v", err)
	}
}
