package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

type Connection struct {
	DB *sql.DB
}

func NewConnection() *Connection {
	conn, err := sql.Open("postgres", getDSN())
	if err != nil {
		slog.Error("failed to establish database connection", "error", err)
	}
	return &Connection{DB: conn}
}

func (conn *Connection) Close() {
	_ = conn.DB.Close()
	conn.DB = nil
}

func getDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
}
