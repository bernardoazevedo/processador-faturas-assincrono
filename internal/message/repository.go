package message

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bernardoazevedo/faturas/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AdicionaNotificacao(queueName string, messageBody []byte) error {
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

func RetornaDelivery(queueName string) (<-chan amqp.Delivery, error) {
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

func EnviaNotificacoes() error {
	horaAtual := utils.RetornaHoraMinutoSegundo()
	utils.WriteLog("\n->started listening for messages: " + horaAtual + "\n")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	amqpMessages, err := RetornaDelivery("notifications")
	if err != nil {
		return err
	}

	go func() {
		for message := range amqpMessages {
			bytes, err := utils.WriteLog(fmt.Sprintf("send: %v", string(message.Body)))
			if err != nil {
				log.Println("error: " + err.Error())
			} else {
				log.Printf("write: %v", bytes)
			}
		}
	}()

	log.Println("[*] Monitoring messages. Press CTRL+C to exit")
	<-sigchan

	defer AMQPconn.Close()

	log.Println("Killed, shutting down")

	return nil
}
