package viewpoint

const (
	AllJobHours = `
		SELECT CONCAT(TRIM(PRTH.Job), ' - ', JCJM.Description) AS Job, 
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
	JobList = `
		SELECT JCJM.Job AS Job, JCJM.Description AS Description
		FROM JCJM 
		JOIN JCCM ON JCJM.JCCo = JCCM.JCCo AND JCJM.Contract = JCCM.Contract
		WHERE JCJM.JCCo=1
		AND JCCM.Department IN(:jobcostdepts)
		AND JCJM.JobStatus=1
	`
)
