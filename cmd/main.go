package main

import (
	"kaffein/config"
	"kaffein/server"
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
	//router := gin.Default()
	//routes.SetupRoute(router)

	// Inisialisasi server
	srv := &server.Application{
		Server: configServer,
		PgSQL:  configPgsql,
	}

	srv.ServeHTTP()
}
