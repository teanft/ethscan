package controller

import (
	"github.com/INFURA/go-ethlibs/eth"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/config"
	"github.com/teanft/ethscan/model"
	"github.com/teanft/ethscan/util"
	"math"
	"math/big"
	"strconv"
	"sync"
)

func SignHandler(c *gin.Context) {
	var block model.EVM
	err := c.Bind(&block)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	signedTxData, err := common.GetEIP155SignedTxData(c, block.Gas, block.Private, block.To, block.Value)
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
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	tx := new(types.Transaction)
	err = common.SendSignedTransaction(c, block.SignedTx, tx)
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
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	rawTx, err := common.GetEIP155RawTx(c, block.Gas, block.Private, block.To, block.Value)
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
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	tx := new(types.Transaction)
	err = common.SendRawedTransaction(c, block.RawTx, tx)
	if err != nil {
		common.Fail(c, gin.H{"send Transaction failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"txHash": tx.Hash().Hex()}, "Success")
}

func TransactionHandler(c *gin.Context) {
	var tx model.Transaction
	err := c.Bind(&tx)
	if err != nil {
		common.Fail(c, gin.H{"bind JSON failed": err.Error()}, "Fail")
		return
	}

	if !util.IsValidHash(tx.TxHash) {
		common.Fail(c, gin.H{"txHash is not valid": tx.TxHash}, "Fail")
		return
	}

	tx, err = GetTransactionData(c, tx.TxHash)
	if err != nil {
		common.Fail(c, gin.H{"get transaction data failed": err.Error()}, "Fail")
		return
	}

	common.Success(c, gin.H{"transaction": tx}, "Success")
}

func TxsHandler(c *gin.Context) {
	blockParam := c.Query("block")

	blockNumber, err := strconv.Atoi(blockParam)
	if err != nil {
		common.Fail(c, gin.H{"block param is not valid": err.Error()}, "Fail")
		return
	}

	page := 1
	if p := c.Query("p"); p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			common.Fail(c, gin.H{"page param is not valid": err.Error()}, "Fail")
			return
		}

		if page <= 0 {
			page = 1
		}
	}

	data, err := common.GetBlockByNumber(c, big.NewInt(int64(blockNumber)))
	if err != nil {
		common.Fail(c, gin.H{"get block failed": err.Error()}, "Fail")
		return
	}

	blockInfo := data.(*types.Block)
	txsLen := len(blockInfo.Transactions())
	pageSize := config.Cfg.Client.PageSize
	maxPage := int(math.Ceil(float64(txsLen) / float64(pageSize)))

	if page >= maxPage {
		page = maxPage
	}

	txs := make([]model.Transaction, txsLen)
	var wg sync.WaitGroup
	startIndex := (page - 1) * pageSize
	endIndex := int(math.Min(float64(txsLen), float64(page*pageSize)))

	for index, transaction := range blockInfo.Transactions() {
		if index >= startIndex && index < endIndex {
			wg.Add(1)
			go func(tx *types.Transaction, txsPtr *[]model.Transaction, index int) {
				defer wg.Done()
				//if *transaction.To() == ethcommon.HexToAddress("0x0000000000000000000000000000000000000000") {
				//    // 忽略合约创建交易
				//    continue
				//}

				txData, err := GetTransactionData(c, tx.Hash().Hex())
				if err != nil {
					common.Fail(c, gin.H{"get transaction data failed": err.Error()}, "Fail")
					return
				}
				(*txsPtr)[index-startIndex] = txData
				//*txsPtr = append(*txsPtr, txData)
			}(transaction, &txs, index)
		}
	}

	wg.Wait()
	txs = txs[:endIndex-startIndex] // 去掉前面的0值

	common.Success(c, gin.H{"transactions": txs, "page": page}, "Success")
}

func GetTransactionData(c *gin.Context, hash string) (model.Transaction, error) {
	tx := model.Transaction{}
	txData, isPending, err := common.GetTransactionByHash(c, hash)
	if err != nil {
		return tx, util.NewErr("get transaction failed", err)
	}

	if isPending {
		return tx, util.NewErr("pending transaction", err)
	}

	txInfo := txData.(*eth.Transaction)

	receiptData, err := common.GetReceiptByHash(c, hash)
	if err != nil {
		return tx, util.NewErr("get receipt failed", err)
	}

	receiptInfo := receiptData.(*eth.TransactionReceipt)
	status := receiptInfo.Status.String()
	switch status {
	case "0x0":
		status = "Failure"
	case "0x1":
		status = "Success"
	}

	data, err := common.GetBlockByNumber(c, receiptInfo.BlockNumber.Big())
	if err != nil {
		return tx, util.NewErr("get block failed", err)
	}

	blockInfo := data.(*types.Block)
	timestamp := blockInfo.Time()
	fee := new(big.Int).Mul(receiptInfo.GasUsed.Big(), txInfo.GasPrice.Big())
	burnt := new(big.Int).Mul(receiptInfo.GasUsed.Big(), blockInfo.BaseFee())
	var toAddress string
	if receiptInfo.To == nil {
		toAddress = receiptInfo.ContractAddress.String()
	} else {
		toAddress = receiptInfo.To.String()
	}

	txSavingsFees := new(big.Int).Sub(txInfo.MaxFeePerGas.Big(), new(big.Int).Add(blockInfo.BaseFee(), txInfo.MaxPriorityFeePerGas.Big()))
	switch txSavingsFees.Cmp(new(big.Int)) {
	case -1, 0:
		txSavingsFees = new(big.Int)
	case 1:
		txSavingsFees = new(big.Int).Mul(txSavingsFees, receiptInfo.GasUsed.Big())
	}

	tx = model.Transaction{
		TxHash:               txInfo.Hash.String(),
		Status:               status,
		BlockNumber:          receiptInfo.BlockNumber.Big(),
		BlockHash:            receiptInfo.BlockHash.String(),
		Timestamp:            timestamp,
		From:                 receiptInfo.From.String(),
		To:                   toAddress,
		Value:                txInfo.Value.Big(),
		TxFee:                fee,
		GasPrice:             txInfo.GasPrice.Big(),
		GasLimit:             txInfo.Gas.UInt64(),
		GasUsed:              receiptInfo.GasUsed.UInt64(),
		BaseFeePerGas:        blockInfo.BaseFee(),
		MaxFeePerGas:         txInfo.MaxFeePerGas.Big(),
		MaxPriorityFeePerGas: txInfo.MaxPriorityFeePerGas.Big(),
		Burnt:                burnt,
		TxnSavingsFees:       txSavingsFees,
		TxnType:              uint8(txInfo.Type.UInt64()),
		Nonce:                txInfo.Nonce.UInt64(),
		TransactionIndex:     uint(receiptInfo.TransactionIndex.UInt64()),
		InputData:            txInfo.Input.String(),
		Logs:                 receiptInfo.Logs,
	}

	return tx, nil
}
