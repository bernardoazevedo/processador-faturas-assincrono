package fatura

type Fatura struct {
	Id         string  `json:"id" bson:"id"`
	Cnpj       string  `json:"cnpj" bson:"cnpj"`
	ValorTotal float64 `json:"valorTotal" bson:"valorTotal"`
	Descricao  string  `json:"descricao" bson:"descricao"`
}
