package util

import (
	"fmt"
	"testing"
)

func TestWalletGenerater(t *testing.T) {
	privateKey, publicKey, address, err := NewWallet()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("privateKey:", privateKey, "\npublicKey:", publicKey, "\naddress:", address)
}
