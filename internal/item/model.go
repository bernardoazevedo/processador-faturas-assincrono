package item

type Item struct {
	Id 				string 	`json:"id"`
	Descricao 		string	`json:"descricao"`
	PrecoCentavos 	int 	`json:"preco_centavos"`
}