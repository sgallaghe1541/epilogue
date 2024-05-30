package viewpoint

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

func ConnectToViewpoint() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	vpServer := os.Getenv("VP_SERVER")
	vpDB := os.Getenv("VP_DATABASE")
	vpUser := os.Getenv("VP_USER")
	vpPass := os.Getenv("VP_PASSWORD")
	vpPort := os.Getenv("VP_PORT")
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%s", vpServer, vpDB, vpUser, vpPass, vpPort)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, fmt.Errorf("viewpoint connection failed: %s", err.Error())
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
