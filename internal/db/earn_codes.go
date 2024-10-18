package db

import "fmt"

type EarnCode struct {
	Description string `db:"description"`
}

func (e EarnCode) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Description, e.Description)
}

func (e EarnCode) SelectString() string {
	return string(e.Description)
}

func (e *EpilogueConnection) GetEarnCodes() ([]EarnCode, error) {
	query := `
		SELECT * FROM earn_codes
	`
	earnCodes := []EarnCode{}

	err := e.DB.Select(&earnCodes, query)
	if err != nil {
		return earnCodes, err
	}
	return earnCodes, nil
}
