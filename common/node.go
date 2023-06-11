package common

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/INFURA/go-ethlibs/node"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/util"
	"math/big"
)

type EvmClientImpl struct {
	ethClient  *ethclient.Client
	libsClient node.Client
}

func NewEVMClient(url string) (*EvmClientImpl, error) {
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		return nil, util.NewErr("failed to create ethclient", err)
	}

	libsClient, err := node.NewClient(context.Background(), url)
	if err != nil {
		return nil, util.NewErr("failed to create libsClient", err)
	}

	return &EvmClientImpl{ethClient: ethClient, libsClient: libsClient}, nil
}

func getClient(c *gin.Context) (EvmClient, error) {
	client, ok := c.Get("ethClient")
	if !ok {
		return nil, errors.New("get ethclient instance fail")
	}

	return client.(EvmClient), nil
}

func GetBlockNumber(ctx *gin.Context) (uint64, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	number, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, util.NewErr("failed to get block number", err)
	}

	return number, nil
}

func GetGasPrice(ctx *gin.Context) (*big.Int, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, util.NewErr("failed to get gas price", err)
	}

	return gasPrice, nil
}

func GetBalanceAt(ctx *gin.Context, address string, blockNumber *big.Int) (*big.Int, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	balance, err := client.BalanceAt(context.Background(), address, blockNumber)
	if err != nil {
		return nil, util.NewErr("failed to get balance", err)
	}

	return balance, nil
}

func GetNonceAt(ctx *gin.Context, address string, blockNumber *big.Int) (uint64, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	nonce, err := client.NonceAt(context.Background(), address, blockNumber)
	if err != nil {
		return 0, util.NewErr("failed to get nonce", err)
	}

	return nonce, nil
}

func GetPendingNonceAt(ctx *gin.Context, address string) (uint64, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return 0, util.NewErr("failed to get pending nonce", err)
	}

	return nonce, err
}

func EstimateGas(ctx *gin.Context, msg any) (uint64, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	gas, err := client.EstimateGas(ctx, msg.(ethereum.CallMsg))
	if err != nil {
		return 0, util.NewErr("failed to estimate gas", err)
	}

	return gas, nil
}

func GetBlockByNumber(ctx *gin.Context, number *big.Int) (any, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	data, err := client.BlockByNumber(context.Background(), number)
	if err != nil {
		return nil, util.NewErr("failed to get block by number", err)
	}

	return data, nil
}

func GetEIP155RawTx(ctx *gin.Context, gas uint64, private, to string, value *big.Int) (string, error) {
	signedTX, err := GetEIP155SignedTx(ctx, gas, private, to, value)
	rawTxBytes, err := signedTX.MarshalBinary()
	if err != nil {
		return "", util.NewErr("cannot marshalBinary", err)
	}

	rawTxHex := hex.EncodeToString(rawTxBytes)

	return rawTxHex, nil
}

func GetEIP155SignedTxData(ctx *gin.Context, gas uint64, private, to string, value *big.Int) (string, error) {
	tx, err := GetEIP155SignedTx(ctx, gas, private, to, value)
	data, err := tx.MarshalBinary()
	if err != nil {
		return "", util.NewErr("cannot marshalBinary", err)
	}

	return hexutil.Encode(data), nil
}

func GetEIP155SignedTx(ctx *gin.Context, gas uint64, private, to string, value *big.Int) (*types.Transaction, error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	toAddress := common.HexToAddress(to)
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return nil, util.NewErr("cannot hexToECDSA", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress.Hex())
	if err != nil {
		return nil, util.NewErr("cannot get PendingNonceAt", err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Gas:   gas,
		Value: value,
	}
	gasLimit, err := EstimateGas(ctx, msg)
	if err != nil {
		return nil, util.NewErr("cannot get EstimateGas", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, util.NewErr("cannot get SuggestGasPrice", err)
	}

	var data []byte

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, util.NewErr("cannot get NetworkID", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, util.NewErr("cannot sign Tx", err)
	}

	return signedTx, nil
}

func SendSignedTransaction(ctx *gin.Context, sign string, tx *types.Transaction) error {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	bytes, err := hexutil.Decode(sign)
	if err != nil {
		return util.NewErr("failed to decode sign", err)
	}

	if err = rlp.DecodeBytes(bytes, &tx); err != nil {
		return util.NewErr("failed to decode transaction", err)
	}

	if err = client.SendTransaction(ctx, tx); err != nil {
		return util.NewErr("failed to send transaction", err)
	}

	return nil
}

func SendRawedTransaction(ctx *gin.Context, rawTx string, tx *types.Transaction) error {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return util.NewErr("failed to decode raw transaction", err)
	}

	if err = rlp.DecodeBytes(rawTxBytes, &tx); err != nil {
		return util.NewErr("failed to decode transaction", err)
	}

	if err = client.SendTransaction(ctx, tx); err != nil {
		return util.NewErr("failed to send transaction", err)
	}

	return nil
}

func GetTransactionByHash(ctx *gin.Context, txHash string) (tx any, isPending bool, err error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	data, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, isPending, util.NewErr("failed to get transaction by txHash", err)
	}

	return data, isPending, nil
}

func GetReceiptByHash(ctx *gin.Context, txHash string) (receipt any, err error) {
	client, err := getClient(ctx)
	if err != nil {
		Fail(ctx, gin.H{"get ethClient error": err}, "Fail")
	}

	data, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, util.NewErr("failed to get receipt by hash", err)
	}

	return data, nil
}
