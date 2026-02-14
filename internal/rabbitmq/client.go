package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

func NewClient() *Client {
	base := os.Getenv("RABBITMQ_HTTP_URL")
	if base == "" {
		base = "http://localhost:15672"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	pass := os.Getenv("RABBITMQ_PASSWORD")
	if pass == "" {
		pass = "guest"
	}

	return &Client{
		baseURL:  strings.TrimRight(base, "/"),
		username: user,
		password: pass,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) DeclareQueue(queue string) error {
	payload := map[string]any{
		"auto_delete": false,
		"durable":     true,
		"arguments":   map[string]any{},
	}

	return c.doRequest(http.MethodPut, "/api/queues/%2F/"+url.PathEscape(queue), payload, nil)
}

func (c *Client) Publish(queue string, rawJSON []byte) error {
	payload := map[string]any{
		"properties":       map[string]any{},
		"routing_key":      queue,
		"payload":          string(rawJSON),
		"payload_encoding": "string",
	}

	return c.doRequest(http.MethodPost, "/api/exchanges/%2F/amq.default/publish", payload, nil)
}

func (c *Client) GetOne(queue string) ([]byte, error) {
	payload := map[string]any{
		"count":    1,
		"ackmode":  "ack_requeue_false",
		"encoding": "auto",
		"truncate": 50000,
	}

	var out []map[string]any
	if err := c.doRequest(http.MethodPost, "/api/queues/%2F/"+url.PathEscape(queue)+"/get", payload, &out); err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, nil
	}

	payloadValue, ok := out[0]["payload"].(string)
	if !ok {
		return nil, fmt.Errorf("response payload was not found")
	}

	return []byte(payloadValue), nil
}

func (c *Client) doRequest(method, path string, body any, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("rabbitmq returned status %d: %s", res.StatusCode, string(respBody))
	}

	if out != nil {
		return json.NewDecoder(res.Body).Decode(out)
	}

	return nil
}
