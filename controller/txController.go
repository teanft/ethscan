package controller

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/model"
)

func SignHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	signedTxData, err := common.GetSignedTxData(block.Gas, block.Private, block.To, block.Value)
	if err != nil {
		common.Fail(c, gin.H{"sign Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"signed_tx_data": signedTxData}, "Success")
}

func SendSignTransactionHandler(c *gin.Context) {
	var block model.Block
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

	common.Success(c, gin.H{"tx_hash": tx.Hash().Hex()}, "Success")
}

func RawHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	rawTx, err := common.GetRawTx(block.Gas, block.Private, block.To, block.Value)
	if err != nil {
		common.Fail(c, gin.H{"raw Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"raw_tx": rawTx}, "Success")
}

func SendRawTransactionHandler(c *gin.Context) {
	var block model.Block
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
