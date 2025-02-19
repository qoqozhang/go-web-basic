package application

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/qoqozhang/go-web-basic/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Application struct {
	DB          *sql.DB
	Engine      *gin.Engine
	Logger      *logger.Helper
	Config      map[string]any
	Middlewares []any
	server      *http.Server
}

func (app *Application) Init() error {
	s := &http.Server{
		Addr:                         "0.0.0.0:8080",
		Handler:                      app.Engine,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	app.server = s
	return nil
}
func (app *Application) Start() error {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := app.Stop(); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	app.server = &http.Server{
		Addr:                         "0.0.0.0:8080",
		Handler:                      app.Engine,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	return nil
}
func (app *Application) Stop() error {
	return app.server.Shutdown(context.Background())
}

func (app *Application) GracefulStop() error {
	return nil
}

// Close 执行数据关闭和清理的工作
func (app *Application) Close() {

}
