package common

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/teanft/ethscan/config"
	"github.com/teanft/ethscan/util"
	"math/big"
)

var Client *ethclient.Client

func NewClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Cfg.Client.URL)
	if err != nil {
		return nil, util.NewErr("failed to create ethclient", err)
	}
	Client = client
	return Client, nil
}

func GetBlockNumber() (uint64, error) {
	number, err := Client.BlockNumber(context.Background())
	if err != nil {
		return 0, util.NewErr("failed to get block number", err)
	}

	return number, nil
}

func GetGasPrice() (*big.Int, error) {
	gasPrice, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, util.NewErr("failed to get gas price", err)
	}

	return gasPrice, nil
}

func GetBalanceAt(address string, blockNumber *big.Int) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := Client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		return nil, util.NewErr("failed to get balance", err)
	}

	return balance, nil
}

func GetPendingBalanceAt(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	pendingBalance, err := Client.PendingBalanceAt(context.Background(), account)
	if err != nil {
		return nil, util.NewErr("failed to get pending balance", err)
	}

	return pendingBalance, nil
}

func GetNonceAt(address string, blockNumber *big.Int) (uint64, error) {
	account := common.HexToAddress(address)
	nonce, err := Client.NonceAt(context.Background(), account, blockNumber)
	if err != nil {
		return 0, util.NewErr("failed to get nonce", err)
	}

	return nonce, nil
}

func GetPendingNonceAt(address string) (uint64, error) {
	account := common.HexToAddress(address)
	nonce, err := Client.PendingNonceAt(context.Background(), account)
	if err != nil {
		return 0, util.NewErr("failed to get pending nonce", err)
	}

	return nonce, err
}

func GetTransactionCountByHash(hash string) (uint, error) {
	blockHash := common.HexToHash(hash)
	count, err := Client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return 0, util.NewErr("failed to get block transaction count", err)
	}

	return count, nil
}

func GetPendingTransactionCount() (uint, error) {
	count, err := Client.PendingTransactionCount(context.Background())
	if err != nil {
		return 0, util.NewErr("failed to get pending transaction count", err)
	}

	return count, nil
}

func SendTransaction(rawTx string, tx *types.Transaction) error {
	rawTxBytes, err := hex.DecodeString(rawTx)

	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return util.NewErr("failed to decode transaction", err)
	}
	err = Client.SendTransaction(context.Background(), tx)
	if err != nil {
		return util.NewErr("failed to send transaction", err)
	}

	return nil
}

func EstimateGas(msg ethereum.CallMsg) (uint64, error) {
	gas, err := Client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, util.NewErr("failed to estimate gas", err)
	}

	return gas, nil
}

func GetBlockByHash(hash string) (*types.Block, error) {
	block, err := Client.BlockByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, util.NewErr("failed to get block by hash", err)
	}

	return block, nil
}

func GetBlockByNumber(number *big.Int) (*types.Block, error) {
	block, err := Client.BlockByNumber(context.Background(), number)
	if err != nil {
		return nil, util.NewErr("failed to get block by number", err)
	}

	return block, nil
}

func GetTransactionByHash(hash string) (*types.Transaction, bool, error) {
	tx, isPending, err := Client.TransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, isPending, util.NewErr("failed to get transaction by hash", err)
	}

	return tx, isPending, nil
}

func GetTransactionSender(tx *types.Transaction, block string, index uint) (common.Address, error) {
	blockHash := common.HexToHash(block)
	address, err := Client.TransactionSender(context.Background(), tx, blockHash, index)
	if err != nil {
		return common.Address{}, util.NewErr("failed to get transaction sender", err)
	}

	return address, nil
}

func GetTransactionInBlock(blockHash string, index uint) (*types.Transaction, error) {
	tx, err := Client.TransactionInBlock(context.Background(), common.HexToHash(blockHash), index)
	if err != nil {
		return nil, util.NewErr("failed to get transaction in block", err)
	}

	return tx, nil
}

func GetTransactionReceipt(hash string) (*types.Receipt, error) {
	receipt, err := Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, util.NewErr("failed to get transaction receipt", err)
	}

	return receipt, nil
}
