package db

import "fmt"

type Employee struct {
	Employee   string `db:"employee"`
	Name       string `db:"name"`
	Class      string `db:"class"`
	Salaried   int64  `db:"salaried"`
	Department string `db:"department"`
}

func (e Employee) SelectValue() string {
	return e.Employee
}

func (e Employee) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (e *EpilogueConnection) GetEmployeesByDepartment(dept string) ([]Employee, error) {
	query := `
		SELECT employee, name, class, salaried, department
		FROM employees
		WHERE active = 1
		AND department = ?
	`
	emps := []Employee{}

	err := e.DB.Select(&emps, query, dept)
	if err != nil {
		return emps, err
	}
	return emps, nil
}

func (e *EpilogueConnection) GetEmployeeByNumber(employeeNumber string) (Employee, error) {
	query := `
		SELECT employee, name, class, salaried, department
		FROM employees
		WHERE active = 1
		AND employee = ?
	`
	emp := Employee{}

	err := e.DB.Select(&emp, query, employeeNumber)
	if err != nil {
		return emp, err
	}
	return emp, nil
}
