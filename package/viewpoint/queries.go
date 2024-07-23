package viewpoint

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	JobListQuery = `
		SELECT Job AS job, Description AS description
		FROM JCJM 
		WHERE JCCo=1 
		AND Job LIKE :jobending 
		AND JobStatus=1 
		AND udFMTS='Y'
	`
	JobHours = `
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
	EmployeeHours = `
		WITH employeehours AS (
			SELECT STR(PRTH.Employee) AS employee, CONCAT(PREH.FirstName, ' ', PREH.LastName) AS name, PRTH.Hours AS hours
			FROM PRTH
			JOIN PREH ON PRTH.PRCo = PREH.PRCo AND PRTH.Employee = PREH.Employee
			WHERE PRTH.PRCo = 1 
			AND PRTH.PRGroup <> 2
			AND PRTH.PRDept IN (:payrolldepts)
			AND PRTH.PREndDate = :wedate
			AND PRTH.Hours <> 0)
		SELECT employee, name, SUM(hours) AS hours
		FROM employeehours
		GROUP BY employee, name
		UNION ALL 
		SELECT 'Total', ' - ', SUM(hours) AS hours
		FROM employeehours
		ORDER BY employee, name
	`
)

type Excelable interface {
	Headers() []interface{}
	DataArray() []interface{}
}

type Job struct {
	Job         string         `db:"job"`
	Description sql.NullString `db:"description"`
}

type JobHoursResult struct {
	Job           sql.NullString  `db:"job"`
	Description   sql.NullString  `db:"description"`
	EEHours       float32         `db:"emphours"`
	EQHours       sql.NullFloat64 `db:"equiphours"`
	EEHoursDetail []EmployeeHoursResult
	EQHoursDetail []EquipmentHoursResult
}

type EmployeeHoursResult struct {
	Employee string  `db:"employee"`
	Name     string  `db:"name"`
	Hours    float32 `db:"hours"`
}

type EquipmentHoursResult struct {
	Equipment   string  `db:"equipment"`
	Description string  `db:"description"`
	Hours       float32 `db:"hours"`
}

func (job *JobHoursResult) GetDetail(args QueryArgs, vpconn *sqlx.DB) {
	job.EEHoursDetail = []EmployeeHoursResult{}
	eerows, err := vpconn.Queryx(JobEmployeeHours, job.Job.String, args.StartWEDate, args.EndWEDate)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = sqlx.StructScan(eerows, &job.EEHoursDetail)
	if err != nil {
		fmt.Print(err.Error())
	}
	eerows.Close()

	eqrows, err := vpconn.Queryx(JobEquipmentHours, job.Job.String, args.StartWEDate, args.EndWEDate)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = sqlx.StructScan(eqrows, &job.EQHoursDetail)
	if err != nil {
		fmt.Print(err.Error())
	}
	eqrows.Close()
}
