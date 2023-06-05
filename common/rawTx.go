package common

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/teanft/ethscan/util"
	"log"
	"math/big"
)

func GetEIP155RawTx(gas uint64, private, to string, value *big.Int) (string, error) {
	signedTX, err := GetEIP155SignedTx(gas, private, to, value)
	rawTxBytes, err := signedTX.MarshalBinary()
	if err != nil {
		return "", util.NewErr("cannot marshalBinary", err)
	}

	rawTxHex := hex.EncodeToString(rawTxBytes)

	return rawTxHex, nil
}

func GetEIP155SignedTxData(gas uint64, private, to string, value *big.Int) (string, error) {
	tx, err := GetEIP155SignedTx(gas, private, to, value)
	data, err := tx.MarshalBinary()
	if err != nil {
		return "", util.NewErr("cannot marshalBinary", err)
	}

	return hexutil.Encode(data), nil
}

func GetEIP155SignedTx(gas uint64, private, to string, value *big.Int) (*types.Transaction, error) {
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
	nonce, err := Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, util.NewErr("cannot get PendingNonceAt", err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Gas:   gas,
		Value: value,
	}
	gasLimit, err := EstimateGas(msg)
	if err != nil {
		return nil, util.NewErr("cannot get EstimateGas", err)
	}
	//gasLimit := uint64(21000) // in units
	gasPrice, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, util.NewErr("cannot get SuggestGasPrice", err)
	}

	var data []byte

	chainID, err := Client.NetworkID(context.Background())
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

func GetEIP1559SignedTx(gas uint64, private, to string, value *big.Int) (*types.Transaction, error) {
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
	nonce, err := Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, util.NewErr("cannot get PendingNonceAt", err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Gas:   gas,
		Value: value,
	}

	gasLimit, err := EstimateGas(msg)
	if err != nil {
		return nil, util.NewErr("cannot get EstimateGas", err)
	}

	tip, _ := GetSuggestGasTipCap()   // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(20000000000) // maxFeePerGas = 20 Gwei
	if err != nil {
		log.Fatal(err)
	}

	var data []byte

	chainID, err := Client.NetworkID(context.Background())
	if err != nil {
		return nil, util.NewErr("cannot get NetworkID", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return nil, util.NewErr("cannot sign Tx", err)
	}
	fmt.Println(signedTx)

	err = Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, util.NewErr("cannot send Tx", err)
	}

	fmt.Printf("Transaction hash: %s", signedTx.Hash().Hex())

	return signedTx, nil
}
