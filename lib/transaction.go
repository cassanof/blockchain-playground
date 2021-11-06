package lib

type Transaction struct {
	Amount float64 `json:"amount"`
	Payer  string  `json:"payer"`
	Payee  string  `json:"payee"`
}
