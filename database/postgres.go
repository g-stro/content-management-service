package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Connection struct {
	DB *sql.DB
}

func NewConnection() (*Connection, error) {
	conn, err := sql.Open("postgres", getDSN())
	if err != nil {
		return nil, err
	}
	return &Connection{DB: conn}, nil
}

func (conn *Connection) Close() {
	_ = conn.DB.Close()
	conn.DB = nil
}

func getDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s timezone=%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL_MODE"), os.Getenv("DB_TIMEZONE"))
}
