package app

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"golang.org/x/oauth2"
)

func NewServer(
	auth *oauth2.Config,
	logger *slog.Logger,
	viewpoint *sqlx.DB,
	epilogue *sqlx.DB,
	users *db.UserModel,
	refreshTokens *db.RefreshTokenModel,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		auth,
		logger,
		viewpoint,
		users,
		refreshTokens,
	)

	var handler http.Handler = mux

	// handler = middleware.Authenticate(logger, handler)

	return handler
}
