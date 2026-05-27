package natsclient

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(ctx context.Context, url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats: connect: %w", err)
	}

	return &NatsClient{conn: conn}, nil
}

func (c *NatsClient) Close() error {
	if err := c.conn.Drain(); err != nil {
		return fmt.Errorf("nats: close: %w", err)
	}

	c.conn.Close()

	return nil
}

func (c *NatsClient) Client() *nats.Conn {
	return c.conn
}
