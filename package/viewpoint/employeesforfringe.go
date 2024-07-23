package viewpoint

import "time"

const (
	EmployeesForFringe = `
		SELECT 
			STR(Employee) AS employee, 
			CONCAT(FirstName, ' ', LastName) AS name,
			CASE 
				WHEN HireDate > COALESCE(RecentRehireDate, '01/01/1900') THEN HireDate
				ELSE RecentRehireDate
			END AS hiredate, 
			HrlyRate AS payrate
		FROM PREH
		WHERE PRCo = 1
		AND ActiveYN = 'Y'
		AND Craft = '1'
		AND PRDept IN (:payrolldepts)
	`
)

type EmployeesForFringeResult struct {
	Employee string    `db:"employee"`
	Name     string    `db:"name"`
	HireDate time.Time `db:"hiredate"`
	Rate     float32   `db:"payrate"`
}

func (e *EmployeesForFringeResult) Headers() *[]interface{} {
	headers := make([]interface{}, 4)
	headers[0] = "Employee"
	headers[1] = "Name"
	headers[2] = "Hire Date"
	headers[3] = "Pay Rate"
	return &headers
}

func (e *EmployeesForFringeResult) DataArray() *[]interface{} {
	data := make([]interface{}, 4)
	data[0] = e.Employee
	data[1] = e.Name
	data[2] = e.HireDate.Format("01/02/2006")
	data[3] = e.Rate
	return &data
}
