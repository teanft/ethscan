package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/util"
	"log"
)

func BlockNumberHandler(c *gin.Context) {
	method := "eth_blockNumber"
	params := make([]string, 0)
	responseData, err := common.CallRPC(c, method, params)
	if err != nil {
		log.Fatal(err)
	}

	result, ok := responseData["result"].(string)
	if !ok {
		common.Fail(c, gin.H{}, "Failed to get result from response")
		return
	}

	blockNumber, err := util.ParseHexStrToInt(result)

	common.Success(c, gin.H{"blockNumber": blockNumber}, "Success")
}
