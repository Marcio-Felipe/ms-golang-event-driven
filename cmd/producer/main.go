package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/example/ms-golang-event-driven/internal/rabbitmq"
)

const queueName = "order.events"

type placeOrderRequest struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
}

type orderPlacedEvent struct {
	EventID   string            `json:"event_id"`
	EventType string            `json:"event_type"`
	Source    string            `json:"source"`
	Occurred  time.Time         `json:"occurred_at"`
	Data      placeOrderRequest `json:"data"`
}

func main() {
	client := rabbitmq.NewClient()
	if err := client.DeclareQueue(queueName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req placeOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		if err := validate(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event := orderPlacedEvent{
			EventID:   req.OrderID,
			EventType: "order.placed",
			Source:    "order-service",
			Occurred:  time.Now().UTC(),
			Data:      req,
		}

		body, err := json.Marshal(event)
		if err != nil {
			http.Error(w, "failed to serialize event", http.StatusInternalServerError)
			return
		}

		if err := client.Publish(queueName, body); err != nil {
			http.Error(w, "failed to publish order event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":   "accepted",
			"message":  "order accepted and published for payment processing",
			"order_id": req.OrderID,
		})
	})

	log.Println("order-service running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}

func validate(req placeOrderRequest) error {
	if strings.TrimSpace(req.OrderID) == "" {
		return errors.New("order_id is required")
	}
	if strings.TrimSpace(req.CustomerID) == "" {
		return errors.New("customer_id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	return nil
}
