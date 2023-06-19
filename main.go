package main

import (
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/config"
	"github.com/teanft/ethscan/middleware"
	"github.com/teanft/ethscan/route"
	"log"
	"os"
)

func start(host, port string) {
	r := gin.Default()
	r.Use(middleware.PanicHandler)
	r = route.CollectRoute(r)

	if port != "" {
		if err := r.Run(host + ":" + port); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	workDir, err := os.Getwd()
	path := workDir + "/config"
	cfg, err := config.InitConfig("application", "yaml", path)
	if err != nil {
		panic(err)
	}

	host := cfg.Server.Host
	port := cfg.Server.Port

	start(host, port)
}
