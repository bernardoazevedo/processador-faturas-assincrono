package fatura

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bernardoazevedo/faturas/internal/logger"
	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/paemuri/brdoc"
)

func QueueList(faturas []Fatura) error {

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

// Simulando chamada para API externa
func generateNote(fatura Fatura) error {
	time.Sleep(time.Second)
	_, err := logger.Add("Nota fiscal emitida para: " + fatura.Id)
	if err != nil {
		return err
	}
	return nil
}
