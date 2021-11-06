package main

import (
	"encoding/json"
	"fmt"

	"github.com/elleven11/blockchain-playground/lib"
)

var satoshi *lib.Wallet
var bob *lib.Wallet
var alice *lib.Wallet

func main() {
	chain := lib.NewChain()

	satoshi = lib.NewWallet(chain)
	bob = lib.NewWallet(chain)
	alice = lib.NewWallet(chain)

	satoshi.SendAmount(50, bob.Public)
	bob.SendAmount(50, bob.Public)
	alice.SendAmount(50, bob.Public)

	// printing blockchain:
	json, err := json.MarshalIndent(chain, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))

}
