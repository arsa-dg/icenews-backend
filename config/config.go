package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func ConnectDB() *sql.DB {
	DBDriver := os.Getenv("DB_DRIVER")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")

	DBUrl := DBDriver + "://" + DBUser + ":" + DBPassword + "@" + DBHost + ":" + DBPort + "/" + DBName

	conCfg, err := pgx.ParseConfig(DBUrl)

	if err != nil {
		log.Fatalln("Parse config failed:", err)
	}

	conCfg.LogLevel = pgx.LogLevelNone

	DB := stdlib.OpenDB(*conCfg)

	return DB
}
