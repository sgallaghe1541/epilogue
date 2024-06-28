package main

import (
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"

	"github.com/sgallaghe1541/epilogue/app"
	"github.com/sgallaghe1541/epilogue/package/viewpoint"
)

// Execute before the service runs.
// func init() {
// 	_ = mime.AddExtensionType(".js", "application/javascript")
// }

func FixMimeTypes() {
	err1 := mime.AddExtensionType(".js", "text/javascript")
	if err1 != nil {
		log.Printf("Error in mime js %s", err1.Error())
	}

	err2 := mime.AddExtensionType(".css", "text/css")
	if err2 != nil {
		log.Printf("Error in mime js %s", err2.Error())
	}
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

	server := &app.Epilogue{
		Logger:    logger,
		Viewpoint: vp,
	}

	server.Logger.Info("starting server")

	err = http.ListenAndServe(addr, server.Routes())
	server.Logger.Error(err.Error())
	os.Exit(1)
}
