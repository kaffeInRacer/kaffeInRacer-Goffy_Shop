package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"kaffein/config"
	"log"
	"sync"
)

var (
	db       *sql.DB
	poolOnce sync.Once
)

func ConnectionPgSQL(c config.PostgresSQL) {
	var err error
	uri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Pass, c.Name, c.SSLMode)
	poolOnce.Do(func() {
		log.Println("[INFO] Connecting to PostgreSQL...")
		defer log.Println("[INFO] Successfully, connection to PostgreSQL closed")
		db, err = sql.Open("pgx", uri)
		if err != nil {
			log.Fatal(err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func GetDBPostgreSQL() *sql.DB {
	return db
}
