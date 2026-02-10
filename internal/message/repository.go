package message

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Add(queueName string, messageBody []byte) error {
	amqpConn, err := GetConn()
	if err != nil {
		return err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        messageBody,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetDelivery(queueName string) (<-chan amqp.Delivery, error) {
	amqpConn, err := GetConn()
	if err != nil {
		return nil, err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return channel.Consume(queue.Name, "", true, false, false, false, nil)
}
