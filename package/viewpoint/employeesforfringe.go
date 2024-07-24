package viewpoint

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	EmployeesForFringe = `
		SELECT 
			TRIM(STR(Employee)) AS employee, 
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

type EmployeeForFringe struct {
	Employee string    `db:"employee"`
	Name     string    `db:"name"`
	HireDate time.Time `db:"hiredate"`
	Rate     float32   `db:"payrate"`
}

func (e *EmployeeForFringe) DataArray() *[]interface{} {
	data := make([]interface{}, 4)
	data[0] = e.Employee
	data[1] = e.Name
	data[2] = e.HireDate.Format("01/02/2006")
	data[3] = fmt.Sprintf("%.2f", e.Rate)
	return &data
}

type EmployeesForFringeResult struct {
	Result []*EmployeeForFringe
}

func (e *EmployeesForFringeResult) Headers() *[]interface{} {
	headers := make([]interface{}, 4)
	headers[0] = "Employee"
	headers[1] = "Name"
	headers[2] = "Hire Date"
	headers[3] = "Pay Rate"
	return &headers
}

func (e *EmployeesForFringeResult) ToExcel(fileName string) error {
	f := excelize.NewFile()

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return err
	}
	err = f.SetSheetRow("Sheet1", "A1", e.Headers())
	if err != nil {
		return err
	}
	err = f.SetCellStyle("Sheet1", "A1", "D1", headerStyle)
	if err != nil {
		return err
	}

	rownum := 2

	for _, emp := range e.Result {
		cell, err := excelize.JoinCellName("A", rownum)
		if err != nil {
			return err
		}
		err = f.SetSheetRow("Sheet1", cell, emp.DataArray())
		if err != nil {
			return err
		}
		rownum++
	}
	f.SaveAs(fileName)
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
