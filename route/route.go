package route

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/block", controller.BlockNumberHandler)
	return r
}
