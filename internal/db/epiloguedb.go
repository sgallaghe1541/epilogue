package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectToEpilogue() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("sqlite3", "internal/db/epilogue.db")
	if err != nil {
		return nil, fmt.Errorf("epilogue connection failed: %s", err.Error())
	}

	return conn, nil
}
