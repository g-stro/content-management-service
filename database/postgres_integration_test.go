//go:build integration

package database

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Fatalf("NewConnection() failed: %v", err)
	}
	defer conn.Close()

	// Test the connection
	err = conn.DB.Ping()
	if err != nil {
		t.Errorf("Database connection failed: %v", err)
	}
}
