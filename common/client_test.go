package common

import (
	"context"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewEVMClient("https://rpc.sepolia.org")
	if err != nil {
		panic(err)
	}

	tx, err := client.libsClient.TransactionByHash(context.Background(), "0x0eca53d26adca5d3a2284283934e07b1a23224b60f15bde9819280109eb177de")
	maxFeePerGasStr := tx.MaxFeePerGas.String()
	maxFeePerGasBig := tx.MaxFeePerGas.Big()

	fmt.Println("maxFeePerGasStr: ", maxFeePerGasStr)
	fmt.Println("maxFeePerGasBig: ", maxFeePerGasBig)
}
