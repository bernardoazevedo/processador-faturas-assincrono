package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	date    time.Time
	message string
}

func getFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteLog(message string) (int, error) {
	date := RetornaDataAtualFormatada()
	fileName := "tmp/" + date + ".txt"

	file, err := getFile(fileName)
	if err != nil {
		log.Println("error opening " + fileName + ": " + err.Error())
		return 0, err
	}
	defer file.Close()

	byteMessage := []byte(message + "\n")
	bytes, err := file.Write(byteMessage)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func formataData(data time.Time) string {
	year := data.Year()
	month := data.Month()
	day := data.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)

}
func RetornaDataAtualFormatada() string {
	return formataData(time.Now())
}

func RetornaHoraMinutoSegundo() string {
	hora, minuto, segundo := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hora, minuto, segundo)
}
