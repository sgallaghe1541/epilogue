package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/db"
)

type App struct {
	Logger    *slog.Logger
	Viewpoint *sqlx.DB
	Epilogue  *sqlx.DB
	Users     *db.UserModel
}
