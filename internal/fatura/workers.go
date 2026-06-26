package fatura

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/dates"
	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/logger"
	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/message"
	"github.com/gin-gonic/gin"
)

func SaveWorker() error {
	horaAtual := dates.ActualDateHMS()
	logger.Add("\t\t\t->started listening for save requests: " + horaAtual)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	amqpMessages, err := message.GetDelivery("save")
	if err != nil {
		return err
	}

	go func() {
		for eachMessage := range amqpMessages {
			var fatura Fatura

			err := json.Unmarshal(eachMessage.Body, &fatura)
			if err != nil {
				log.Println("error parsing note to save: " + err.Error())
			}

			err = save(fatura)
			if err != nil {
				log.Println("error saving note: " + err.Error())
			}

			err = message.Add("generateNote", eachMessage.Body)
			if err != nil {
				log.Println("error adding message to generateNote queue, at item [" + fatura.Id + "]: " + err.Error())
			}
		}
	}()

	log.Println("[*] Monitoring save requests. Press CTRL+C to exit")
	<-sigchan

	log.Println("Killed, shutting down")

	return nil
}

func GenerateNoteWorker() error {
	horaAtual := dates.ActualDateHMS()
	logger.Add("\t\t\t->started listening for note request: " + horaAtual)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	amqpMessages, err := message.GetDelivery("generateNote")
	if err != nil {
		return err
	}

	go func() {
		for eachMessage := range amqpMessages {
			var fatura Fatura

			err := json.Unmarshal(eachMessage.Body, &fatura)
			if err != nil {
				log.Println("error parsing note: " + err.Error())
			}

			err = generateNote(fatura)
			if err != nil {
				log.Println("error creating note: " + err.Error())
			}

			messageBody, err := json.Marshal(gin.H{
				"cnpj":      fatura.Cnpj,
				"descricao": fmt.Sprintf(
					"Foi emitida uma nota fiscal no valor de R$%s com descrição: '%s' no CNPJ: %s", 
					strconv.FormatFloat(fatura.ValorTotal, 'f', -1, 64), fatura.Descricao, fatura.Cnpj,
				),
			})
			if err != nil {
				log.Println("error parsing notification: " + err.Error())
			}

			err = message.Add("notifications", messageBody)
			if err != nil {
				log.Println("error creating notification, at item [" + fatura.Id + "]: " + err.Error())
			}
		}
	}()

	log.Println("[*] Monitoring note requests. Press CTRL+C to exit")
	<-sigchan

	log.Println("Killed, shutting down")

	return nil
}

func NotificationsWorker() error {
	horaAtual := dates.ActualDateHMS()
	logger.Add("\t\t\t->started listening for messages: " + horaAtual)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	amqpMessages, err := message.GetDelivery("notifications")
	if err != nil {
		return err
	}

	go func() {
		for message := range amqpMessages {
			_, err := logger.Add(fmt.Sprintf("send: %v", string(message.Body)))
			if err != nil {
				log.Println("error sending notification: " + err.Error())
			}
		}
	}()

	log.Println("[*] Monitoring messages. Press CTRL+C to exit")
	<-sigchan

	log.Println("Killed, shutting down")

	return nil
}
