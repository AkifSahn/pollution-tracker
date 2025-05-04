package database

import (
	"context"
	"fmt"
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool *pgxpool.Pool
)

func checkSchema() {
	checkTable := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_schema = 'public' AND table_name = 'air_pollution'
        );
    `

	ctx := context.Background()

	_, err := DBPool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS postgis;`)
	if err != nil {
		log.Fatal("Failed to enable PostGIS extension - ", err)
	}

	var exist bool
	err = DBPool.QueryRow(ctx, checkTable).Scan(&exist)
	if err != nil {
		log.Fatal("Failed to check if table exists - ", err)
	}

	if exist {
		log.Println("Table 'air_pollution' already exists")
		return
	}

	log.Println("Table 'air_pollution' does not exists. Creating the table")

	createTable := `
		CREATE TABLE air_pollution (
			time        TIMESTAMPTZ       NOT NULL,
			pollutant   TEXT              NOT NULL,
			value       DOUBLE PRECISION  NOT NULL,
			is_anomaly  BOOLEAN           NOT NULL DEFAULT false,
			latitude    DOUBLE PRECISION  NOT NULL,
			longitude   DOUBLE PRECISION  NOT NULL,
			geog        GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
				ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
			) STORED
		);
    `

	_, err = DBPool.Exec(ctx, createTable)
	if err != nil {
		log.Fatal("Failed to create table - ", err)
	}

	_, err = DBPool.Exec(ctx, `SELECT create_hypertable('air_pollution', 'time');`)
	if err != nil {
		log.Fatal("Failed to create hypertable - ", err)
	}

	log.Println("Table 'air_pollution' and hypertable created successfully.")
}

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
		checkSchema()
	}
}
