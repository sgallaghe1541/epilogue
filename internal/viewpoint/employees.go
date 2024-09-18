package viewpoint

import "fmt"

const (
	employeesByDivision = `
		SELECT Employee AS employee, CONCAT(FirstName, " ", LastName) AS name
		FROM PREH
		WHERE ActiveYN = 'Y'
		AND PRGroup = 1
		AND PRCo = 1
		AND PRDept IN (:payrolldepts)
	`
)

type Employee struct {
	Employee   string `db:"employee"`
	Name       string `db:"name"`
	Class      string `db:"class"`
	Salaried   int64  `db:"salaried"`
	Department string `db:"prdept"`
}

func (e Employee) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (e Employee) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

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
