package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Epilogue struct {
	Logger    *slog.Logger
	Viewpoint *sqlx.DB
}
