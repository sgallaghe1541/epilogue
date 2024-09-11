package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sgallaghe1541/epilogue/internal/auth"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
	"golang.org/x/oauth2"
)

type app struct {
	auth           *oauth2.Config
	logger         *slog.Logger
	viewpoint      *sqlx.DB
	epilogue       *sqlx.DB
	sessionManager *scs.SessionManager
	users          *db.UserModel
	jobs           *viewpoint.TimeEntryModel
	timecards      *db.TimeCardModel
	refreshTokens  *db.RefreshTokenModel
}

func newApp(auth *oauth2.Config,
	logger *slog.Logger,
	viewpoint *sqlx.DB,
	epilogue *sqlx.DB,
	sessionManager *scs.SessionManager,
	users *db.UserModel,
	jobs *viewpoint.TimeEntryModel,
	timecards *db.TimeCardModel,
	refreshTokens *db.RefreshTokenModel) *app {
	return &app{
		auth:           auth,
		logger:         logger,
		viewpoint:      viewpoint,
		epilogue:       epilogue,
		sessionManager: sessionManager,
		users:          users,
		jobs:           jobs,
		timecards:      timecards,
		refreshTokens:  refreshTokens,
	}
}

func run(ctx context.Context) error {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	addr := ":3000"

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	vp, err := viewpoint.ConnectToViewpoint()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	defer vp.Close()

	data, err := db.ConnectToEpilogue()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	defer data.Close()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(data.DB)
	sessionManager.Lifetime = 12 * time.Hour

	app := newApp(
		auth.MicrosoftConfig(),
		logger,
		vp,
		data,
		sessionManager,
		&db.UserModel{DB: data},
		&viewpoint.TimeEntryModel{DB: vp},
		&db.TimeCardModel{DB: data},
		&db.RefreshTokenModel{DB: data},
	)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	logger.Info(fmt.Sprintf("starting server. listening on %s", httpServer.Addr))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
