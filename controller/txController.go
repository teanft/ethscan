package controller

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	common "github.com/teanft/ethscan/common"

	"github.com/teanft/ethscan/model"
)

func SignHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	signedTxData, err := common.GetEIP155SignedTxData(block.Gas, block.Private, block.To, block.Value)
	if err != nil {
		common.Fail(c, gin.H{"sign Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"signedTxData": signedTxData}, "Success")
}

func SendSignTransactionHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	tx := new(types.Transaction)
	err = common.SendSignedTransaction(block.SignedTx, tx)
	if err != nil {
		common.Fail(c, gin.H{"send Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"txHash": tx.Hash().Hex()}, "Success")
}

func RawHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	rawTx, err := common.GetEIP155RawTx(block.Gas, block.Private, block.To, block.Value)
	if err != nil {
		common.Fail(c, gin.H{"raw Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"rawTx": rawTx}, "Success")
}

func SendRawTransactionHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	tx := new(types.Transaction)
	err = common.SendRawedTransaction(block.RawTx, tx)
	if err != nil {
		common.Fail(c, gin.H{"sendTransaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"txHash": tx.Hash().Hex()}, "Success")
}

func TxHandler(c *gin.Context) {
	var tx model.Transaction
	err := c.Bind(&tx)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	txInfo, isPending, err := common.GetTransactionByHash(tx.TxHash)
	if err != nil {
		common.Fail(c, gin.H{"get transaction failed": err.Error()}, "Fail")
		return
	}

	if isPending {
		common.Fail(c, gin.H{"transaction is pending": err.Error()}, "Fail")
		return
	}

	tx = model.Transaction{
		//TransactionIndex: txInfo,
		TxHash: txInfo.Hash().Hex(),
		//BlockNumber: txInfo,
		//Timestamp: ,
		//From: txInfo,
		To:       txInfo.To().Hex(),
		Value:    txInfo.Value(),
		GasUsed:  txInfo.Gas(),
		GasPrice: txInfo.GasPrice(),
		//BaseFeePerGas: txInfo,
		MaxFeePerGas:         txInfo.GasFeeCap(),
		MaxPriorityFeePerGas: txInfo.GasTipCap(),
		//BlockHash: ,
		InputData: string(txInfo.Data()),
		Nonce:     txInfo.Nonce(),
		TxnType:   txInfo.Type(),
		Size:      txInfo.Size(),
		Cost:      txInfo.Cost(),
	}
}
