package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type App struct {
	Logger    *slog.Logger
	Viewpoint *sqlx.DB
	Epilogue  *sqlx.DB
}
