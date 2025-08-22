package message

import (
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
)

var AMQPconn *amqp.Connection

func Connect() (*amqp.Connection, error) {
	user := os.Getenv("RABBIT_USER")
	pass := os.Getenv("RABBIT_PASS")

	conn, err := amqp.Dial("amqp://"+user+":"+pass+"@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	AMQPconn = conn
	return AMQPconn, nil
}

func Start() (error) {
	_, err := Connect()
	if err != nil {
		return err
	}
	return nil
}

func GetConn() (*amqp.Connection) {
	return AMQPconn
}