# Event-Driven Order and Payment Services with RabbitMQ

A professional, minimal Go example of asynchronous communication between microservices using RabbitMQ.

## Architecture

This repository contains two services connected by RabbitMQ:

1. **order-service**
   - Exposes an HTTP API to place orders.
   - Publishes an `order.placed` event to the `order.events` queue.

2. **payment-service**
   - Consumes `order.placed` events from the `order.events` queue.
   - Simulates payment processing and logs payment outcome.

RabbitMQ acts as the event bridge between these services.

## Tech Stack

- Go 1.22
- RabbitMQ (management image)
- Docker Compose

## Run with Docker Compose

```bash
docker compose up --build
```

### Endpoints and Ports

- Order Service API: `http://localhost:8080`
- RabbitMQ Management UI: `http://localhost:15672`
- RabbitMQ credentials: `guest` / `guest`

## Place an Order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORD-1001",
    "customer_id": "CUS-789",
    "amount": 149.90,
    "currency": "USD"
  }'
```

Expected response:

```json
{
  "message": "order accepted and published for payment processing",
  "order_id": "ORD-1001",
  "status": "accepted"
}
```

## Check Payment Processing

Watch the `payment-service` logs:

```bash
docker compose logs -f payment-service
```

Expected log format:

```text
[PAYMENT_PROCESSED] order_id=ORD-1001 customer_id=CUS-789 amount=149.90 currency=USD payment_status=approved duration_ms=512
```

## Run Locally Without Docker (Optional)

1. Start RabbitMQ with management plugin enabled.
2. Run services in separate terminals:

```bash
RABBITMQ_HTTP_URL=http://localhost:15672 go run ./cmd/consumer
```

```bash
RABBITMQ_HTTP_URL=http://localhost:15672 go run ./cmd/producer
```

3. Call `POST /orders` as shown above.
