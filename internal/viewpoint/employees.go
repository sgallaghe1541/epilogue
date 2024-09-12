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
	Employee string `db:"employee"`
	Name     string `db:"name"`
}

func (e Employee) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (e Employee) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (v *ViewpointConnection) GetEmployeesByDivision(div QueryArgs) ([]Employee, error) {
	emps := []Employee{}

	query, args, err := BuildInQuery(employeesByDivision, div)
	if err != nil {
		return nil, err
	}
	err = v.DB.Select(&emps, query, args)
	if err != nil {
		return nil, err
	}
	return emps, nil
}
