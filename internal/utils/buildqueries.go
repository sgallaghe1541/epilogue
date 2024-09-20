package utils

import (
	"fmt"

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
