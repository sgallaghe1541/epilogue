package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/sgallaghe1541/epilogue/app"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

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

	data, err := db.ConnectToEpilogue()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer data.Close()

	server := &app.App{
		Logger:    logger,
		Viewpoint: vp,
		Epilogue:  data,
	}

	server.Logger.Info("starting server")

	err = http.ListenAndServe(addr, server.Routes())
	server.Logger.Error(err.Error())
	os.Exit(1)
}
