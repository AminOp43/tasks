package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitDB() *sql.DB {
	connStr := "user=postgres dbname=test_db password=12345678 sslmode=disable"
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
