package fatura

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bernardoazevedo/faturas/internal/database"
	"github.com/bernardoazevedo/faturas/internal/dates"
	"github.com/bernardoazevedo/faturas/internal/logger"
	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/paemuri/brdoc"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProcessList(faturas []Fatura) error {

	for _, fatura := range faturas {
		faturaJson, err := json.Marshal(fatura)
		if err != nil {
			return err
		}

		err = validate(fatura)
		if err != nil {
			return err
		}

		err = message.Add("save", faturaJson)
		if err != nil {
			return errors.New("error adding message to save queue, at item: " + fatura.Id)
		}

		
		err = message.Add("generateNote", faturaJson)
		if err != nil {
			return errors.New("error adding message to generate queue, at item: " + fatura.Id)
		}

		messageBody, err := json.Marshal(gin.H{
			"cnpj":      fatura.Cnpj,
			"descricao": fmt.Sprintf("Foi emitida uma nota fiscal no valor de R$%s com descrição: '%s' no CNPJ: %s", strconv.FormatFloat(fatura.ValorTotal, 'f', -1, 64), fatura.Descricao, fatura.Cnpj),
		})
		if err != nil {
			return err
		}

		err = message.Add("notifications", messageBody)
		if err != nil {
			return errors.New("error creating notification at item: " + fatura.Id)
		}
	}

	return nil
}

func save(fatura Fatura) (error) {
	DB := database.GetDB()
	
	collection  := DB.Database("faturasAPI").Collection("faturas")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, fatura)
	if err != nil {
		return errors.New("error inserting item: " + fatura.Id)
	}

	return nil
}

func validate(fatura Fatura) error {
	if fatura.ValorTotal <= 0 {
		return errors.New("the total value must be above 0, at item: " + fatura.Id)
	}

	if strings.TrimSpace(fatura.Descricao) == "" {
		return errors.New("the description can't be empty, at item: " + fatura.Id)
	}

	if !validateCnpj(fatura.Cnpj) {
		return errors.New("the cnpj is invalid, at item: " + fatura.Id)
	}

	return nil
}

func validateCnpj(cnpj string) bool {
	return brdoc.IsCNPJ(cnpj)
}

func List() ([]Fatura, error) {
	DB := database.GetDB()

	collection  := DB.Database("faturasAPI").Collection("faturas")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("error searching item")
	}
	defer cursor.Close(ctx)

	var faturas []Fatura
	for cursor.Next(ctx) {
		var fatura Fatura
		if err := cursor.Decode(&fatura); err != nil {
			return nil, errors.New("error decoding item")
		}
		faturas = append(faturas, fatura)
	}

	return faturas, nil
}

// Simulando chamada para API externa
func generateNote(fatura Fatura) error {
	time.Sleep(time.Second)
	_, err := logger.Add("Nota fiscal emitida para: " + fatura.Id)
	if err != nil {
		return err
	}
	return nil
}

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
		for message := range amqpMessages {
			var fatura Fatura

			err := json.Unmarshal(message.Body, &fatura)
			if err != nil {
				log.Println("error: " + err.Error())
			}

			err = save(fatura)
			if err != nil {
				log.Println("error: " + err.Error())
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
		for message := range amqpMessages {
			var fatura Fatura

			err := json.Unmarshal(message.Body, &fatura)
			if err != nil {
				log.Println("error: " + err.Error())
			}

			err = generateNote(fatura)
			if err != nil {
				log.Println("error: " + err.Error())
			}
		}
	}()

	log.Println("[*] Monitoring note requests. Press CTRL+C to exit")
	<-sigchan

	log.Println("Killed, shutting down")

	return nil	
}