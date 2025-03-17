package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
)

type configurationApp interface {
	ServerApp() ServerApp
	PostgresSQL() PostgresSQL
}

type confApp struct {
	serverApp   ServerApp
	postgresSQL PostgresSQL
}

func LoadConfig(path string) configurationApp {
	var once sync.Once
	var conf *confApp
	once.Do(func() {
		if err := godotenv.Load(path); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		conf = &confApp{
			serverApp: ServerApp{
				Host: getEnvStr("SERVER_HOST", "localhost"),
				Port: getEnvStr("SERVER_PORT", "8080"),
			},
			postgresSQL: PostgresSQL{
				Host:    getEnvStr("POSTGRES_HOSTNAME", "localhost"),
				User:    getEnvStr("POSTGRES_USER", "postgres"),
				Pass:    getEnvStr("POSTGRES_PASSWORD", ""),
				Port:    getEnvStr("POSTGRES_PORT", "5432"),
				Name:    getEnvStr("POSTGRES_DB", "postgres"),
				SSLMode: getEnvStr("POSTGRES_SSL_MODE", "disable"),
			},
		}
	})
	return conf
}

type ServerApp struct {
	Host string
	Port string
}

func (c *confApp) ServerApp() ServerApp {
	return c.serverApp
}

type PostgresSQL struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

func (c *confApp) PostgresSQL() PostgresSQL {
	return c.postgresSQL
}

func getEnvStr(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultVal
}
