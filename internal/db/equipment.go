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

func (e *EpilogueConnection) GetEquipmentByDepartment(dept string) ([]Equipment, error) {
	query := `
		SELECT equipment, description, department
		FROM equipment
		WHERE active = 1
		AND department = ?
	`
	equipment := []Equipment{}

	err := e.DB.Select(&equipment, query, dept)
	if err != nil {
		return equipment, err
	}
	return equipment, nil
}

func (e *EpilogueConnection) GetEquipmentByNumber(equipmentNumber string) (Equipment, error) {
	query := `
		SELECT equipment, description, department
		FROM equipment
		WHERE active = 1
		AND equipment = ?
	`
	equipment := Equipment{}

	err := e.DB.Select(equipment, query, equipmentNumber)
	if err != nil {
		return equipment, err
	}
	return equipment, nil
}
