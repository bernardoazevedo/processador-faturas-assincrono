package message

import (
	"errors"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var AMQPconn *amqp.Connection

func Connect() (*amqp.Connection, error) {
	user := os.Getenv("RABBIT_USER")
	pass := os.Getenv("RABBIT_PASS")

	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open test channel: %w", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare("connection.test", false, true, true, false, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("connection verification failed: %w", err)
	}

	channel.QueueDelete("connection.test", false, false, false)

	AMQPconn = conn
	return AMQPconn, nil
}

func Start() error {
	var lastErr error

	for range 5 {
		conn, err := Connect()
		if err == nil && conn != nil {
			return nil
		}
		lastErr = err
		time.Sleep(2 * time.Second)
	}

	return errors.New("failed to connect to RabbitMQ after 5 attempts: " + lastErr.Error())
}

func GetConn() (*amqp.Connection, error) {
	if AMQPconn.IsClosed() {
		return Connect()
	}
	return AMQPconn, nil
}
