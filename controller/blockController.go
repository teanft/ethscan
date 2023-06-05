package controller

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/model"
	"math/big"
)

func BlockHeight(c *gin.Context) {
	blockNumber, err := common.GetBlockNumber()
	if err != nil {
		common.Fail(c, gin.H{"get LastBlockNumber failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"blockNumber": blockNumber}, "Success")
}

func GasPriceHandler(c *gin.Context) {
	gasPrice, err := common.GetGasPrice()
	if err != nil {
		common.Fail(c, gin.H{"get GasPrice failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"gasPrice": gasPrice}, "Success")
}

func BalanceHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
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
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	nonce, err := common.GetNonceAt(block.Address, block.BlockNumber)
	if err != nil {
		common.Fail(c, gin.H{"get Nonce failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"nonce": nonce}, "Success")
}

func PendingNonceHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	nonce, err := common.GetPendingNonceAt(block.Address)
	if err != nil {
		common.Fail(c, gin.H{"get pending nonce failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"nonce": nonce}, "Success")
}

func BlockHandler(c *gin.Context) {
	var block model.Block
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	blockInfo, err := common.GetBlockByNumber(block.Number)
	if err != nil {
		common.Fail(c, gin.H{"get block failed": err.Error()}, "Fail")
		return
	}

	burntFees := big.NewInt(0)
	burntFees.Mul(big.NewInt(int64(blockInfo.GasUsed())), blockInfo.BaseFee())

	block = model.Block{
		Number:        blockInfo.Number(),
		Timestamp:     blockInfo.Time(),
		Miner:         blockInfo.Coinbase().Hex(),
		Size:          blockInfo.Size(),
		GasUsed:       blockInfo.GasUsed(),
		GasLimit:      blockInfo.GasLimit(),
		BaseFeePerGas: blockInfo.BaseFee(),
		BurntFees:     burntFees,
		ExtraData:     hexutil.Encode(blockInfo.Extra()[:]),
		//ExtraData:       string(blockInfo.Extra()),
		Hash:            blockInfo.Hash().Hex(),
		ParentHash:      blockInfo.ParentHash().Hex(),
		StateRoot:       blockInfo.Header().Root.Hex(),
		WithdrawalsRoot: blockInfo.Header().WithdrawalsHash.Hex(),
		Nonce:           blockInfo.Nonce(),
		Difficulty:      blockInfo.Difficulty(),
	}

	common.Success(c, gin.H{"block": block}, "Success")
}
