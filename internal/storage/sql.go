package storage

import (
	"context"
	"database/sql"
	"time"
)

type Config struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func ConnectSQL(c Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.Dsn)
	if err != nil {
		return nil, err
	}

	parseMaxIdleTime, err := time.ParseDuration(c.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(parseMaxIdleTime)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
