package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qoqozhang/go-web-basic/application"
	"github.com/qoqozhang/go-web-basic/logger"
	"github.com/qoqozhang/go-web-basic/logger/zap"
)

func main() {
	g := gin.Default()
	zaplogger, _ := zap.NewZapLogger(logger.SetLevel(logger.DebugLevel))
	log := logger.NewHelper(zaplogger)
	app := new(application.Application)
	app.Logger = log
	app.Engine = g

}
