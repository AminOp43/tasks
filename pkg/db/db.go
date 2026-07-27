package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() *sql.DB {
	connStr := "user=postgres dbname=messages_db password=12345678 sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return nil
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database not responding:", err)
		return nil
	}
	fmt.Println("Connected to database successfully!")
	return DB
}
