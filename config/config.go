package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

func ConnectDB() *pgx.Conn {
	DBDriver := os.Getenv("DB_DRIVER")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")

	DBUrl := DBDriver + "://" + DBUser + ":" + DBPassword + "@" + DBHost + ":" + DBPort + "/" + DBName

	connConfig, err := pgx.ParseConfig(DBUrl)

	if err != nil {
		log.Fatalln(err)
	}

	DB, err := pgx.ConnectConfig(context.Background(), connConfig)

	if err != nil {
		log.Fatalln(err)
	}

	return DB
}
