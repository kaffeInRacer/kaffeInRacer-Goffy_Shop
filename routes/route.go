package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoute(router *gin.Engine) *gin.Engine {
	// Contoh route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	router.POST("/post", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	})

	router.GET("/ping", func(c *gin.Context) {
		x := []rune("abcdefghijklmnopq")
		fmt.Println(x[:1000])
		c.JSON(200, gin.H{"message": "pong"})
	})

	return router
}
