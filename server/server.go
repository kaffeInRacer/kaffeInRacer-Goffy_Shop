package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"kaffein/config"
	"kaffein/pkg/database/postgresql"
	"kaffein/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Server config.ServerApp
	PgSQL  config.PostgresSQL
}

func (app *Application) ServeHTTP() {
	addr := fmt.Sprintf("%s:%s", app.Server.Host, app.Server.Port)
	gin.SetMode(gin.DebugMode)
	postgresql.ConnectionPgSQL(app.PgSQL)

	srv := http.Server{
		Addr:    addr,
		Handler: routes.SetupRoute(gin.Default()),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Printf("[INFO] Start server at %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	<-quit
	log.Println("[WARNING] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ALERT] Server Shutdown Failed: %+v", err)
	}

	log.Println("[INFO] Server gracefully stopped")
}
