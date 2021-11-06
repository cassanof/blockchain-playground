package lib

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
)

type SignPair struct {
	Signature []byte
	HashSum   []byte
}

type Wallet struct {
	Public     string `json:"public"`
	Private    string `json:"private"`
	RSAPrivate *rsa.PrivateKey
	RSAPublic  *rsa.PublicKey
	Chain      *Chain
}

func NewWallet(chain *Chain) *Wallet {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	public := &private.PublicKey

	privateBytes, err := x509.MarshalPKCS8PrivateKey(private)
	publicBytes, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		panic(err)
	}

	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	}

	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	}

	privateEncoded := pem.EncodeToMemory(privateBlock)
	publicEncoded := pem.EncodeToMemory(publicBlock)
	if err != nil {
		panic(err)
	}

	wallet := Wallet{
		Private:    string(privateEncoded),
		Public:     string(publicEncoded),
		RSAPrivate: private,
		RSAPublic:  public,
		Chain:      chain,
	}
	return &wallet
}

func (w *Wallet) SendAmount(amount float64, payeePublic string) {
	trans := Transaction{amount, w.Public, payeePublic}
	transEncoded, err := json.Marshal(trans)
	if err != nil {
		panic(err)
	}

	// signing transaction
	hasher := crypto.SHA256.New()
	hasher.Write(transEncoded)

	hashSum := hasher.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, w.RSAPrivate, crypto.SHA256, hashSum, nil)
	if err != nil {
		panic(err)
	}

	signPair := SignPair{
		Signature: signature,
		HashSum:   hashSum,
	}

	w.Chain.addBlock(trans, *w.RSAPublic, signPair)
}
