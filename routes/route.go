package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"kaffein/modules/users/usersRepository"
	"kaffein/pkg/database/postgresql"
	"log"
)

func SetupRoute(router *gin.Engine) *gin.Engine {

	db := postgresql.GetDBPostgreSQL()
	//user := &[]domain.Users{}

	UserRepository := usersRepository.NewUsersRepository(db)
	res, err := UserRepository.FetchAll(context.Background())
	if err != nil {
		log.Fatal("Error fetching user:", err)
	}

	// Cek apakah user ditemukan sebelum mencetak
	for _, re := range res {
		fmt.Println(re)
	}

	return router
}
