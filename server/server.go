package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"kaffein/config"
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

func (app *Application) ServeHTTP(router *gin.Engine) {
	addr := fmt.Sprintf("%s:%s", app.Server.Host, app.Server.Port)

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Printf("Server running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}

	log.Println("Server gracefully stopped")
}
