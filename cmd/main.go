package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

type application struct {
	logger    *slog.Logger
	viewpoint *sql.DB
}

func main() {
	addr := ":4000"
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	vp, err := viewpoint.ConnectToViewpoint()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer vp.Close()

	app := &application{
		logger:    logger,
		viewpoint: vp,
	}

	app.logger.Info("starting server")

	err = http.ListenAndServe(addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
