package route

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/height", controller.BlockHeight)
	r.GET("/gas_price", controller.GasPriceHandler)
	r.POST("/balance", controller.BalanceHandler)
	r.POST("/nonce", controller.NonceHandler)
	r.POST("/pending_nonce", controller.PendingNonceHandler)
	return r
}
