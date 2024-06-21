package viewpoint

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func BuildInQuery(querystring string, div Division) (string, []interface{}, error) {
	query, args, err := sqlx.Named(querystring, div)
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
