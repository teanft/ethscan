package common

import (
	"context"
	"github.com/INFURA/go-ethlibs/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/teanft/ethscan/util"
	"math/big"
)

func (c *EvmClientImpl) BlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

func (c *EvmClientImpl) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.ethClient.SuggestGasPrice(ctx)
}

func (c *EvmClientImpl) BalanceAt(ctx context.Context, address string, blockNumber *big.Int) (*big.Int, error) {
	return c.ethClient.BalanceAt(ctx, common.HexToAddress(address), blockNumber)
}

func (c *EvmClientImpl) NonceAt(ctx context.Context, address string, blockNumber *big.Int) (uint64, error) {
	//account := common.HexToAddress(address)
	//return c.ethClient.NonceAt(ctx, account, blockNumber)
	account, _ := eth.NewAddress(address)
	numArg := util.ToBlockNumArg(blockNumber)
	blockNumberOrTag, _ := eth.NewBlockNumberOrTag(numArg)
	return c.libsClient.GetTransactionCount(ctx, *account, *blockNumberOrTag)
}
func (c *EvmClientImpl) PendingNonceAt(ctx context.Context, address string) (uint64, error) {
	return c.ethClient.PendingNonceAt(ctx, common.HexToAddress(address))
}

func (c *EvmClientImpl) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.ethClient.NetworkID(ctx)
}

func (c *EvmClientImpl) EstimateGas(ctx context.Context, msg any) (uint64, error) {
	return c.ethClient.EstimateGas(ctx, msg.(ethereum.CallMsg))
}

func (c *EvmClientImpl) SendTransaction(ctx context.Context, tx any) error {
	return c.ethClient.SendTransaction(ctx, tx.(*types.Transaction))
}

func (c *EvmClientImpl) BlockByNumber(ctx context.Context, number *big.Int) (interface{}, error) {
	return c.ethClient.BlockByNumber(ctx, number)
}

func (c *EvmClientImpl) TransactionByHash(ctx context.Context, txHash string) (tx any, isPending bool, err error) {
	transaction, err := c.libsClient.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, false, err
	}

	return transaction, transaction.BlockNumber.Big() == nil, err
}

func (c *EvmClientImpl) TransactionReceipt(ctx context.Context, txHash string) (any, error) {
	return c.libsClient.TransactionReceipt(ctx, txHash)
}
