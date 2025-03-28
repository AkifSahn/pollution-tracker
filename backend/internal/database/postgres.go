package database

import (
	"context"
	"fmt"
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/jackc/pgx/v5"
)

var (
	DB *pgx.Conn
)

func connectDB(cfg *config.Config) *pgx.Conn {
	ctx := context.Background()
	// postgres://username:password@host:port/dbname
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	//run a simple query to check our connection
	var greeting string
	err = conn.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)
	if err != nil {
		log.Fatal("QueryRow failed: ", err)
	}

	return conn
}

func InitDB(cfg *config.Config) {
	if DB == nil {
		db := connectDB(cfg)
		DB = db
	}
}

func NewTransaction(ctx context.Context) (pgx.Tx, error) {
	return DB.Begin(ctx)
}
