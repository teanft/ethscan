package common

import (
	"context"
	"math/big"
)

type EvmClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account string, blockNumber *big.Int) (*big.Int, error)
	NonceAt(ctx context.Context, account string, blockNumber *big.Int) (uint64, error)
	PendingNonceAt(ctx context.Context, account string) (uint64, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg any) (uint64, error)
	SendTransaction(ctx context.Context, tx any) error
	BlockByNumber(ctx context.Context, number *big.Int) (any, error)
	TransactionByHash(ctx context.Context, txHash string) (tx any, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash string) (any, error)
}
