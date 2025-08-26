package fatura

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bernardoazevedo/faturas/internal/database"
	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/paemuri/brdoc"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ProcessaFaturas(faturas []Fatura) error {

	for _, fatura := range faturas {
		err := validaFatura(fatura)
		if err != nil {
			return err
		}

		err = salvaFatura(fatura)
		if err != nil {
			return err
		}

		messageBody, err := json.Marshal(gin.H{
			"cnpj":      fatura.Cnpj,
			"descricao": fmt.Sprintf("Foi emitida uma nota fiscal no valor de R$%s com descrição: '%s' no CNPJ: %s", strconv.FormatFloat(fatura.ValorTotal, 'f', -1, 64), fatura.Descricao, fatura.Cnpj),
		})
		if err != nil {
			return err
		}

		err = message.AdicionaNotificacao("notifications", messageBody)
		if err != nil {
			return errors.New("error creating notification at item: " + fatura.Id)
		}
	}

	return nil
}

func salvaFatura(fatura Fatura) (error) {
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

func validaFatura(fatura Fatura) error {
	if fatura.ValorTotal <= 0 {
		return errors.New("the total value must be above 0, at item: " + fatura.Id)
	}

	if strings.TrimSpace(fatura.Descricao) == "" {
		return errors.New("the description can't be empty, at item: " + fatura.Id)
	}

	if !validaCnpj(fatura.Cnpj) {
		return errors.New("the cnpj is invalid, at item: " + fatura.Id)
	}

	return nil
}

func validaCnpj(cnpj string) bool {
	return brdoc.IsCNPJ(cnpj)
}

func ListaFaturas() ([]Fatura, error) {
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
