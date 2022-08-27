package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func DatabaseConnection() (*sql.DB, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("not DATABASE .ENV FILE")
	}

	database_connection := os.Getenv("DATABASE_CONNECTION")

	url := database_connection

	res, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}

	if _, err := res.Exec("CREATE TABLE IF NOT EXISTS book (name TEXT, author TEXT, id int NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))"); err != nil {
		return res, nil
	}

	if err := res.Ping(); err != nil {
		return nil, err
	}

	return res, nil
}
