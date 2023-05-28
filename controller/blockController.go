package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/model"
)

func BlockHeight(c *gin.Context) {
	blockNumber, err := common.GetBlockNumber()
	if err != nil {
		common.Fail(c, gin.H{"getLastBlockNumber failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"blockNumber": blockNumber}, "Success")
}

func GasPriceHandler(c *gin.Context) {
	gasPrice, err := common.GetGasPrice()
	if err != nil {
		common.Fail(c, gin.H{"getGasPrice failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"gasPrice": gasPrice}, "Success")
}

func BalanceHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	// 是否包含合约地址
	if block.Contract != "" {
		// TODO: 调用合约进行查询
	}

	balance, err := common.GetBalanceAt(block.Address, block.BlockNumber)
	if err != nil {
		common.Fail(c, gin.H{"getBalanceAt failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"balance": balance}, "Success")
}

func NonceHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	nonce, err := common.GetNonceAt(block.Address, block.BlockNumber)
	if err != nil {
		common.Fail(c, gin.H{"getNonceAt failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"nonce": nonce}, "Success")
}

func PendingNonceHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bindJSON failed": err.Error()}, "Fail")
		return
	}

	nonce, err := common.GetPendingNonceAt(block.Address)
	if err != nil {
		common.Fail(c, gin.H{"getNonceAt failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"nonce": nonce}, "Success")
}
