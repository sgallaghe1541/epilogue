package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type EpilogueConnection struct {
	DB *sqlx.DB
}

func ConnectToEpilogue() (*EpilogueConnection, error) {
	conn, err := sqlx.Connect("sqlite3", "internal/db/epilogue.db")
	if err != nil {
		return nil, fmt.Errorf("epilogue connection failed: %s", err.Error())
	}

	return &EpilogueConnection{DB: conn}, nil
}
