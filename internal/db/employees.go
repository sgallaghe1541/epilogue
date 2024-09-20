package db

import "fmt"

const (
	employeesByDepartment = `
		SELECT employee, name, class, salaried, department
		FROM employees
		WHERE active = 1
		AND department = ?
	`
)

type Employee struct {
	Employee   string `db:"employee"`
	Name       string `db:"name"`
	Class      string `db:"class"`
	Salaried   int64  `db:"salaried"`
	Department string `db:"department"`
}

func (e Employee) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (e Employee) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}
