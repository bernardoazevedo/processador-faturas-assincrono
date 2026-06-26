package main

import (
	"log"

	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/database"
	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/fatura"
	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/message"
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

	go fatura.SaveWorker()
	go fatura.GenerateNoteWorker()
	go fatura.NotificationsWorker()

	router := gin.Default()
	loadRoutes(router)
	router.Run(":1234")
}

func loadRoutes(router *gin.Engine) {
	router.POST("/faturas", fatura.HttpProcessList)
	router.GET("/faturas", fatura.HttpList)
}
