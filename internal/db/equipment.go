package db

import "fmt"

type Equipment struct {
	Equipment   string `db:"equipment"`
	Description string `db:"description"`
	Department  string `db:"department"`
}

func (e Equipment) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Equipment, e.Description)
}

func (e Equipment) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Equipment, e.Description)
}
