package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"kaffein/config"
	"kaffein/routes"
	"kaffein/server"
	"os"
)

var (
	configServer config.ServerApp
	configPgsql  config.PostgresSQL
)

func init() {
	conf := config.LoadConfig(".env")
	configServer = conf.ServerApp()
	configPgsql = conf.PostgresSQL()
}

func main() {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	f, _ = os.Create("gin-error.log")
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)

	router := gin.Default()
	router.Use(gin.ErrorLogger(), gin.Recovery())
	routes.SetupRoute(router)

	srv := &server.Application{
		Server: configServer,
		PgSQL:  configPgsql,
	}

	srv.ServeHTTP(router)
}
