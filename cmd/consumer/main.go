package main

import (
	"encoding/json"
	"log"
	"math/rand"
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

	log.Println("payment-service listening for order events...")
	rand.Seed(time.Now().UnixNano())

	for {
		body, err := client.GetOne(queueName)
		if err != nil {
			log.Printf("failed to fetch event: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if body == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		var event orderPlacedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("invalid event payload: %s", string(body))
			continue
		}

		processPayment(event)
	}
}

func processPayment(event orderPlacedEvent) {
	startedAt := time.Now()

	// Simulated gateway latency
	time.Sleep(time.Duration(300+rand.Intn(400)) * time.Millisecond)

	status := "approved"
	if event.Data.Amount > 10000 {
		status = "manual_review"
	}

	log.Printf(
		"[PAYMENT_PROCESSED] order_id=%s customer_id=%s amount=%.2f currency=%s payment_status=%s duration_ms=%d",
		event.Data.OrderID,
		event.Data.CustomerID,
		event.Data.Amount,
		event.Data.Currency,
		status,
		time.Since(startedAt).Milliseconds(),
	)
}
