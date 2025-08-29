package main

import (
	"log"

	"github.com/bernardoazevedo/faturas/internal/database"
	"github.com/bernardoazevedo/faturas/internal/fatura"
	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetPrefix("main: ")
	log.SetFlags(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.Start()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = message.Start()
	if err != nil {
		log.Fatal("Error connecting to rabbitmq")
	}

	go message.NotificationsWorker()
	go fatura.SaveWorker()
	go fatura.GenerateNoteWorker()

	router := gin.Default()
	loadRoutes(router)
	router.Run(":1234")
}

func loadRoutes(router *gin.Engine) {

	router.POST("/faturas", fatura.HttpProcessaFaturas)
	router.GET("/faturas", fatura.HttpListaFaturas)
}
