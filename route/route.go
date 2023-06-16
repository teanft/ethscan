package route

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/common"
	"github.com/teanft/ethscan/config"
	"github.com/teanft/ethscan/controller"
	"github.com/teanft/ethscan/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	client, err := common.NewEVMClient(config.Cfg.Client.URL)
	if err != nil {
		panic(err)
	}

	r.GET("/height", middleware.WithEVMClient(client), controller.BlockHeightHandler)
	r.GET("/gas_price", middleware.WithEVMClient(client), controller.GasPriceHandler)
	r.POST("/balance", middleware.WithEVMClient(client), controller.BalanceHandler)
	r.POST("/nonce", middleware.WithEVMClient(client), controller.NonceHandler)
	r.POST("/pending_nonce", middleware.WithEVMClient(client), controller.PendingNonceHandler)
	// transaction 组路由
	tx := r.Group("/transaction")
	{
		//tx.POST("", controller.TxHandler)
		tx.POST("/sign", middleware.WithEVMClient(client), controller.SignHandler)
		tx.POST("/raw", middleware.WithEVMClient(client), controller.RawHandler)
		tx.POST("/send_sign", middleware.WithEVMClient(client), controller.SendSignTransactionHandler)
		tx.POST("/send_raw", middleware.WithEVMClient(client), controller.SendRawTransactionHandler)
	}
	r.POST("/tx", middleware.WithEVMClient(client), controller.TransactionHandler)
	r.GET("txs", middleware.WithEVMClient(client), controller.TxsHandler)
	r.POST("/block", middleware.WithEVMClient(client), controller.BlockHandler)
	return r
}
