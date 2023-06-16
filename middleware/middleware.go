package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"net/http"
)

func WithEVMClient(client common.EvmClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ethClient", client)
		c.Next()
	}
}

func PanicHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case *common.ResponseMsg:
				common.Fail(c, gin.H{"panic": t.Msg}, "Fail")
			default:
				common.Response(c, http.StatusOK, http.StatusInternalServerError, gin.H{"panic": "internal error"}, "服务器内部异常")
			}
			c.Abort()
		}
	}()

	c.Next()
}
