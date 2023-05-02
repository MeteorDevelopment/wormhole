package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"wormhole/pkg/config"
)

var dbPool *pgxpool.Pool

func Get() *pgxpool.Pool {
	return dbPool
}

func Init() {
	cfg := config.Get()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	var err error
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping the database: %v", err)
	}
}

func Close() {
	dbPool.Close()
}
