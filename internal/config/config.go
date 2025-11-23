package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

func LoadDBConfig(envFile string) *DBConfig {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("No %s file found, using environment variables", envFile)
		}
	} else {
		_ = godotenv.Load()
	}

	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "10"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifeSec, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "300"))

	return &DBConfig{
		DSN:          getEnv("DB_DSN", ""),
		MaxOpenConns: maxOpen,
		MaxIdleConns: maxIdle,
		ConnMaxLife:  time.Duration(connMaxLifeSec) * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
