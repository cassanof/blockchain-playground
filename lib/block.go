package lib

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Block struct {
	PrevHash    *string      `json:"prevHash"` // Link to previous block
	Transaction *Transaction `json:"transaction"`
	Nonce       float64      `json:"nonce"`
	Ts          time.Time    `json:"time"`
}

func (b *Block) hash() *string {
	json, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(json)

	hash := hex.EncodeToString(hasher.Sum(nil))

	return &hash
}
