package db

type Route struct {
	Path        string       `db:"path"`
	Permissions []Permission `db:"permissions"`
}
