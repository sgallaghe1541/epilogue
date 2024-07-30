package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/sgallaghe1541/epilogue/app"
	"github.com/sgallaghe1541/epilogue/internal/auth"
	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

func main() {
	addr := ":3000"
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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

	sessionManager := scs.New()

	sessionManager.Store = sqlite3store.New(data.DB)

	server := &app.App{
		Logger:         logger,
		Viewpoint:      vp,
		Epilogue:       data,
		Users:          &db.UserModel{DB: data},
		SessionManager: sessionManager,
	}

	auth.NewAuth()

	server.Logger.Info("starting server")

	err = http.ListenAndServe(addr, server.Routes())
	server.Logger.Error(err.Error())
	os.Exit(1)
}
