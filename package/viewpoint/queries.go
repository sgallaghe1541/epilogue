package viewpoint

import "database/sql"

const (
	AllJobHours = `
		SELECT CONCAT(TRIM(PRTH.Job), ' - ', JCJM.Description) AS job, 
			CONCAT(REPLACE(PRTH.Phase, ' ',''), ' - ', JCJP.Description) AS Phase, 
			CONVERT(varchar,PRTH.PostDate,1) AS Date, PRTH.Employee, CONCAT(PREH.FirstName, ' ', PREH.LastName) AS Name, 
			CONCAT(PRTH.EarnCode, ' - ', PREC.Description) AS EarnCode, 
			PRTH.Hours AS Hours
		FROM PRTH
		JOIN PREH ON PRTH.PRCo = PREH.PRCo AND PRTH.Employee = PREH.Employee
		LEFT JOIN JCJM ON PRTH.PRCo = JCJM.JCCo AND PRTH.Job = JCJM.Job
		LEFT JOIN JCJP ON PRTH.JCCo = PRTH.PRCo AND PRTH.Job = JCJP.Job AND PRTH.Phase = JCJP.Phase
		LEFT JOIN PREC ON PRTH.PRCo = PREC.PRCo AND PRTH.EarnCode = PREC.EarnCode
		WHERE PRTH.PRCo = 1 
		AND PRTH.PRGroup <> 2
		AND PRTH.PRDept IN (:payrolldepts)
		AND PRTH.PREndDate = ?
		AND PRTH.Hours <> 0
	`
	JobListQuery = `
		SELECT Job AS job, Description AS description
		FROM JCJM 
		WHERE JCCo=1 
		AND Job LIKE :jobending 
		AND JobStatus=1 
		AND udFMTS='Y'
	`
)

type Job struct {
	Job         string         `db:"job"`
	Description sql.NullString `db:"description"`
}
