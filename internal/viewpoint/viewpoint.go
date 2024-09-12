package viewpoint

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

type ViewpointConnection struct {
	DB *sqlx.DB
}

func ConnectToViewpoint() (*ViewpointConnection, error) {

	vpServer := os.Getenv("VP_SERVER")
	vpDB := os.Getenv("VP_DATABASE")
	vpUser := os.Getenv("VP_USER")
	vpPass := os.Getenv("VP_PASSWORD")
	vpPort := os.Getenv("VP_PORT")
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%s", vpServer, vpDB, vpUser, vpPass, vpPort)

	conn, err := sqlx.Open("mssql", connString)
	if err != nil {
		return nil, fmt.Errorf("viewpoint connection failed: %s", err.Error())
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &ViewpointConnection{DB: conn}, nil
}
