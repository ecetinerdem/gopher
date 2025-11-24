package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ecetinerdem/gopherSocial/internal/env"
)

func New(addr string, maxOpenConns int, maxIdleConns int, maxIdLeTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	duration := env.GetDuration("MAX_IDLE_TIME", maxIdLeTime)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(duration))
	db.SetMaxIdleConns(maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil

}
