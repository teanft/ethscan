package util

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewWallet() (privateKeyString, publicKeyString, address string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", NewErr("failed to generate private key", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString = hexutil.Encode(publicKeyBytes)[4:]

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return privateKeyString, publicKeyString, address, err
}
