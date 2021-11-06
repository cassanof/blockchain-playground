package lib

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Chain struct {
	Blocks *[]Block `json:"blocks"`
}

func NewChain() *Chain {
	return &Chain{
		Blocks: &[]Block{
			{
				PrevHash: nil,
				Transaction: &Transaction{
					Amount: 1337,
					Payer:  "genesis",
					Payee:  "satoshi",
				},
				Ts: time.Now(),
			},
		},
	}
}

func (c *Chain) lastBlock() *Block {
	blk := *c.Blocks
	return &blk[len(blk)-1]
}

func (c *Chain) mine(nonce float64) float64 {
	solution := 1.0
	fmt.Println("mining...")

	for {
		hasher := crypto.MD5.New()
		hasher.Write([]byte(fmt.Sprintf("%f", nonce+solution)))

		attempt := hex.EncodeToString(hasher.Sum(nil))

		if attempt[0:4] == "0000" {
			fmt.Printf("Solved: %f\n", solution)
			return solution
		}
		solution += 1.0
	}
}

func (c *Chain) addBlock(trans Transaction, publicKey rsa.PublicKey, signPair SignPair) {
	err := rsa.VerifyPSS(&publicKey, crypto.SHA256, signPair.HashSum, signPair.Signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}

	newBlock := Block{
		PrevHash:    c.lastBlock().hash(),
		Transaction: &trans,
		Nonce:       math.Round(rand.Float64() * 999999999),
		Ts:          time.Now(),
	}
	c.mine(newBlock.Nonce)
	*c.Blocks = append(*c.Blocks, newBlock)
}
