package main

import (
	"log"

	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/joho/godotenv"
)

func main() {
	log.SetPrefix("workers/message: ")
	log.SetFlags(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = message.Start()
	if err != nil {
		log.Fatal("Error connecting to rabbitmq: " + err.Error())
	}

	message.EnviaNotificacoes()
}
