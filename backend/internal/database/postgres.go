package database

import (
	"context"
	"fmt"
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool *pgxpool.Pool
)

func connectDB(cfg *config.Config) *pgxpool.Pool {
	ctx := context.Background()
	// postgres://username:password@host:port/dbname
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	//run a simple query to check our connection
	var greeting string
	err = pool.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)
	if err != nil {
		log.Fatal("QueryRow failed: ", err)
	}

	return pool
}

func InitDB(cfg *config.Config) {
	if DBPool == nil {
		DBPool = connectDB(cfg)
	}
}

func NewTransaction(ctx context.Context) (pgx.Tx, error) {
	return DBPool.Begin(ctx)
}
