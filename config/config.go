package config

import (
	"context"
	"fmt"
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
	fmt.Println(DBUrl)

	DB, err := pgx.Connect(context.Background(), DBUrl)
	if err != nil {
		fmt.Println("Error while connecting to database!")
	}

	return DB
}
