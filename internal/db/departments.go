package db

type Department struct {
	ID          string `db:"departmentid"`
	Description string `db:"description"`
}
