package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

type DB struct {
	SQL *sql.DB
}

func New() (*DB, error) {
	server := getenv("DB_SERVER", "localhost")
	port := getenv("DB_PORT", "1433")
	user := getenv("DB_USER", "sa")
	pass := os.Getenv("DB_PASSWORD")
	name := getenv("DB_NAME", "mcp_kodepos")

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		user, pass, server, port, name,
	)

	sqlDB, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return &DB{SQL: sqlDB}, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
