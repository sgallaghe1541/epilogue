package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sgallaghe1541/epilogue/app"
	"github.com/sgallaghe1541/epilogue/auth"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
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

	srv := app.NewServer(
		auth.MicrosoftConfig(),
		logger,
		vp,
		data,
		&db.UserModel{DB: data},
		&db.RefreshTokenModel{DB: data},
	)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
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
