# Microserviço simples com RabbitMQ (treino)

Projeto mínimo em Go para praticar comunicação assíncrona com RabbitMQ:

- **producer**: expõe endpoint HTTP e publica mensagens na fila.
- **consumer**: consome a fila e imprime no log.

> Objetivo didático (simples, sem foco em robustez).

## Subir com Docker Compose

```bash
docker compose up --build
```

Serviços:
- Producer HTTP: `http://localhost:8080`
- RabbitMQ UI: `http://localhost:15672` (usuário/senha: `guest`/`guest`)

## Publicar mensagem

```bash
curl -X POST http://localhost:8080/publish \
  -H "Content-Type: application/json" \
  -d '{"name":"joao","message":"ola rabbit"}'
```

Resposta esperada:

```json
{"status":"mensagem publicada"}
```

## Ver consumo

Veja os logs do container `consumer`. Ele exibirá algo como:

```text
[EVENTO] name=joao message=ola rabbit created_at=2026-01-01T12:00:00Z
```

## Rodar local sem Docker (opcional)

1. Tenha um RabbitMQ com plugin de management rodando local em `http://localhost:15672`
2. Em terminais separados:

```bash
RABBITMQ_HTTP_URL=http://localhost:15672 go run ./cmd/consumer
```

```bash
RABBITMQ_HTTP_URL=http://localhost:15672 go run ./cmd/producer
```
