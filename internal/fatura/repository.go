package fatura

import (
	"context"
	"errors"
	"time"

	"github.com/bernardoazevedo/processadorFaturasAssincrono/internal/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func save(fatura Fatura) error {
	DB := database.GetDB()

	collection := DB.Database("faturasAPI").Collection("faturas")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, fatura)
	if err != nil {
		return errors.New("error inserting item: " + fatura.Id)
	}

	return nil
}

func List() ([]Fatura, error) {
	DB := database.GetDB()

	collection := DB.Database("faturasAPI").Collection("faturas")
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
