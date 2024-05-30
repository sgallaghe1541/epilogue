package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

type application struct {
	logger    *slog.Logger
	viewpoint *sql.DB
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	vp, err := viewpoint.ConnectToViewpoint()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer vp.Close()

	app := application{
		logger:    logger,
		viewpoint: vp,
	}
}
