package message

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func RetornaNotificacoes(queueName string) (<-chan amqp.Delivery, error) {
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
	date := retornaDataAtualFormatada()
	nomeArquivo := "tmp/" + date + ".txt"

	file, err := os.OpenFile(nomeArquivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error opening " + nomeArquivo + ": " + err.Error())
		return err
	}
	defer file.Close()

	horaAtual := retornaHoraMinutoSegundo()
	bytes, err := file.Write([]byte("\n->started: " + horaAtual + "\n"))
	if err != nil {
		log.Println("error: " + err.Error())
	} else {
		log.Printf("write: %v", bytes)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	amqpMessages, err := RetornaNotificacoes("notifications")
	if err != nil {
		return err
	}

	go func() {
		for message := range amqpMessages {
			byteMessage := []byte(fmt.Sprintf("send: %v\n", string(message.Body)))
			bytes, err := file.Write(byteMessage)
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

func formataData(data time.Time) string {
	year := data.Year()
	month := data.Month()
	day := data.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func retornaDataAtualFormatada() string {
	return formataData(time.Now())
}

func retornaHoraMinutoSegundo() string {
	hora, minuto, segundo := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hora, minuto, segundo)
}
