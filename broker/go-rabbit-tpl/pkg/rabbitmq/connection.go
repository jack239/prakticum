package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Connection{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Connection) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}
