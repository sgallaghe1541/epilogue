package viewpoint

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

const (
	JobHoursQuery = `
	    WITH jobhours AS (
			SELECT PRTH.Job AS job, JCJM.Description AS description, PRTH.Hours AS emphours, PRTH.UsageUnits AS equiphours
			FROM PRTH
			LEFT JOIN JCJM ON PRTH.PRCo = JCJM.JCCo AND PRTH.Job = JCJM.Job
			WHERE PRTH.PRCo = 1 
			AND PRTH.PRGroup <> 2
			AND PRTH.PRDept IN (:payrolldepts)
			AND PRTH.PREndDate BETWEEN :startwedate AND :endwedate)
		SELECT job, description, SUM(emphours) AS emphours, SUM(equiphours) AS equiphours
		FROM jobhours
		GROUP BY job, description
		UNION ALL 
		SELECT 'Total', ' - ', SUM(emphours) AS emphours, SUM(equiphours) AS equiphours
		FROM jobhours
		ORDER BY job, description
	`
	JobEmployeeHours = `
		SELECT STR(PRTH.Employee) AS employee, CONCAT(PREH.FirstName, ' ', PREH.LastName) AS name, SUM(PRTH.Hours) AS hours
		FROM PRTH
		JOIN PREH ON PRTH.PRCo = PREH.PRCo AND PRTH.Employee = PREH.Employee
		WHERE PRTH.PRCo = 1 
		AND PRTH.PRGroup <> 2
		AND PRTH.Job = ?
		AND PRTH.PREndDate BETWEEN ? AND ?
		AND PRTH.Hours <> 0
		GROUP BY PRTH.Employee, PREH.FirstName, PREH.LastName
	`
	JobEquipmentHours = `
		SELECT STR(PRTH.Equipment) AS equipment, EMEM.Description AS description, SUM(PRTH.UsageUnits) AS hours
		FROM PRTH
		JOIN EMEM ON PRTH.PRCo = EMEM.PRCo AND PRTH.Equipment = EMEM.Equipment
		WHERE PRTH.PRCo = 1 
		AND PRTH.PRGroup <> 2
		AND PRTH.Job = ?
		AND PRTH.PREndDate BETWEEN ? AND ?
		AND PRTH.UsageUnits <> 0
		GROUP BY PRTH.Equipment, EMEM.Description
	`
)

type EmployeeHours struct {
	Employee string  `db:"employee"`
	Name     string  `db:"name"`
	Hours    float32 `db:"hours"`
}

func (e *EmployeeHours) DataArray() *[]interface{} {
	data := make([]interface{}, 2)
	data[0] = e.Name
	data[1] = fmt.Sprintf("%.2f", e.Hours)
	return &data
}

type EmployeeHoursResult struct {
	Result []EmployeeHours
}

func (em *EmployeeHoursResult) Headers() *[]interface{} {
	headers := make([]interface{}, 2)
	headers[0] = "Employee"
	headers[1] = "Hours"
	return &headers
}

type EquipmentHours struct {
	Equipment   string  `db:"equipment"`
	Description string  `db:"description"`
	Hours       float32 `db:"hours"`
}

func (e *EquipmentHours) DataArray() *[]interface{} {
	data := make([]interface{}, 2)
	data[0] = fmt.Sprintf("%s--%s", e.Equipment, e.Description)
	data[1] = fmt.Sprintf("%.2f", e.Hours)
	return &data
}

type EquipmentHoursResult struct {
	Result []EquipmentHours
}

func (em *EquipmentHoursResult) Headers() *[]interface{} {
	headers := make([]interface{}, 2)
	headers[0] = "Equipment"
	headers[1] = "Hours"
	return &headers
}

func (job *JobHours) GetDetail(args QueryArgs, vpconn *sqlx.DB) {
	job.EEHoursDetail.Result = []EmployeeHours{}
	eerows, err := vpconn.Queryx(JobEmployeeHours, job.Job.String, args.StartWEDate, args.EndWEDate)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = sqlx.StructScan(eerows, &job.EEHoursDetail.Result)
	if err != nil {
		fmt.Print(err.Error())
	}
	eerows.Close()

	eqrows, err := vpconn.Queryx(JobEquipmentHours, job.Job.String, args.StartWEDate, args.EndWEDate)
	if err != nil {
		fmt.Print(err.Error())
	}
	job.EQHoursDetail.Result = []EquipmentHours{}
	err = sqlx.StructScan(eqrows, &job.EQHoursDetail.Result)
	if err != nil {
		fmt.Print(err.Error())
	}
	eqrows.Close()
}

type JobHours struct {
	Job           sql.NullString  `db:"job"`
	Description   sql.NullString  `db:"description"`
	EEHours       float32         `db:"emphours"`
	EQHours       sql.NullFloat64 `db:"equiphours"`
	EEHoursDetail EmployeeHoursResult
	EQHoursDetail EquipmentHoursResult
}

func (j *JobHours) DataArray() *[]interface{} {
	data := make([]interface{}, 5)
	data[0] = j.Job.String
	data[1] = j.Description.String
	data[2] = ""
	data[3] = fmt.Sprintf("%.2f", j.EEHours)
	data[4] = fmt.Sprintf("%.2f", j.EQHours.Float64)
	return &data
}

type JobHoursResult struct {
	Result []*JobHours
}

func (j *JobHoursResult) Headers() *[]interface{} {
	headers := make([]interface{}, 5)
	headers[0] = "Job"
	headers[1] = "Description"
	headers[2] = ""
	headers[3] = "Employee Hours"
	headers[4] = "Equipment Hours"
	return &headers
}

func (j *JobHoursResult) ToExcel(fileName string) error {
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

	jobStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#9BC2E6"},
		},
	})
	if err != nil {
		return err
	}

	err = f.SetSheetRow("Sheet1", "A1", j.Headers())
	if err != nil {
		return err
	}
	err = f.SetCellStyle("Sheet1", "A1", "E1", headerStyle)
	if err != nil {
		return err
	}
	rowNum := 2
	eecolNum := 1
	eqcolNum := 4

	for _, job := range j.Result {
		cell, err := excelize.JoinCellName("A", rowNum)
		if err != nil {
			return err
		}
		err = f.SetSheetRow("Sheet1", cell, job.DataArray())
		if err != nil {
			return err
		}

		endCell, err := excelize.CoordinatesToCellName(5, rowNum)
		if err != nil {
			return err
		}

		err = f.SetCellStyle("Sheet1", cell, endCell, jobStyle)
		if err != nil {
			return err
		}

		rowNum++
		cell, err = excelize.CoordinatesToCellName(eecolNum, rowNum)
		if err != nil {
			return nil
		}
		err = f.SetSheetRow("Sheet1", cell, job.EEHoursDetail.Headers())
		if err != nil {
			return err
		}

		cell, err = excelize.CoordinatesToCellName(eqcolNum, rowNum)
		if err != nil {
			return err
		}
		err = f.SetSheetRow("Sheet1", cell, job.EQHoursDetail.Headers())
		if err != nil {
			return err
		}
		rowNum++

		for i, ee := range job.EEHoursDetail.Result {
			cell, err = excelize.CoordinatesToCellName(eecolNum, rowNum+i)
			if err != nil {
				return err
			}
			err = f.SetSheetRow("Sheet1", cell, ee.DataArray())
			if err != nil {
				return err
			}
		}
		for i, eq := range job.EQHoursDetail.Result {
			cell, err = excelize.CoordinatesToCellName(eqcolNum, rowNum+i)
			if err != nil {
				return err
			}
			err = f.SetSheetRow("Sheet1", cell, eq.DataArray())
			if err != nil {
				return err
			}
		}
		if len(job.EEHoursDetail.Result) >= len(job.EQHoursDetail.Result) {
			rowNum += len(job.EEHoursDetail.Result) + 1
		} else {
			rowNum += len(job.EQHoursDetail.Result) + 1
		}
		f.SaveAs(fileName)
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
