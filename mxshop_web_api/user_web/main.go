package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop_web_api/initialize"
)

func main() {
	Router := initialize.Routers()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	suger := logger.Sugar()
	port := 8021
	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {

	}
}
