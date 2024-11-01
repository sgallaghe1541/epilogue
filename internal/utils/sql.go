package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Arg interface {
	Arg()
}

func BuildInQuery(querystring string, params Arg) (string, []interface{}, error) {
	query, args, err := sqlx.Named(querystring, params)
	if err != nil {
		return "", nil, fmt.Errorf("failed to prep named query: %s", err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, fmt.Errorf("failed to prep in query: %s", err)
	}

	query = sqlx.Rebind(sqlx.QUESTION, query)
	return query, args, nil
}

func BuildQuery(querystring string, params Arg) (string, []interface{}, error) {
	query, args, err := sqlx.Named(querystring, params)
	if err != nil {
		return "", nil, fmt.Errorf("failed to prep named query: %s", err)
	}
	query = sqlx.Rebind(sqlx.QUESTION, query)
	return query, args, nil
}

func MakeSQLNullString(s string) sql.NullString {
	if s != "" {
		return sql.NullString{String: s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func MakeSQLNullFloat64(f float64) sql.NullFloat64 {
	if f != 0.00 {
		return sql.NullFloat64{Float64: f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

func MakeSQLNullTime(t time.Time) sql.NullTime {
	if !t.IsZero() {
		return sql.NullTime{Time: t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}
