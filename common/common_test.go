package common

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/teanft/ethscan/util"
	"log"
	"math"
	"math/big"
	"testing"
)

//var ec ethclient.Client

func NewTestClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("https://rpc.sepolia.org")
	if err != nil {
		return nil, util.NewErr("failed to create ethclient", err)
	}
	Client = client
	return Client, nil
}

func TestGetBalanceAt(t *testing.T) {
	client, err := NewTestClient()
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0x8bf0f7D37E6f8c8f3B087151825F5CCF23aB1B45")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Println("当前余额为：", balance)

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("ETH的值为：", ethValue)

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("pending中的值为：", pendingBalance)
}

func TestGetRawTx(t *testing.T) {
	_, err := NewTestClient()
	if err != nil {
		log.Fatal(err)
	}
	rawTx, err := getRawTx(21000, "0xc0749b740cAe8768b89547fEdbC33eB45afC236c", big.NewInt(234))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(rawTx)
}

func TestGetSignedTxData(t *testing.T) {
	_, err := NewTestClient()
	if err != nil {
		log.Fatal(err)
	}
	rawTx, err := getSignedTxData(21000, "0xc0749b740cAe8768b89547fEdbC33eB45afC236c", big.NewInt(234))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(rawTx)
}
