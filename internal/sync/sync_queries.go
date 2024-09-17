package sync

const (
	syncAllJobs = `
		SELECT Job AS job, Description AS description, PRStateCode AS state, Certified AS certified, CraftTemplate AS template
		FROM JCJM 
		WHERE JCCo=1 
		AND JobStatus=1 
		AND udFMTS='Y'
	`
	syncAllPhases = `
		SELECT Job AS job, Phase AS phase, Description AS description
		FROM JCJP WHERE udFMTS='Y'
		AND Job IN (
			SELECT Job
			FROM JCJM
			WHERE JCCo=1
			AND JobStatus=1
			AND udFMTS='Y'
		)
	`
	syncAllEmployees = `
		SELECT Employee AS employee,
			CONCAT(FirstName, " ", LastName) AS fullname,
			Class AS class,
			CASE 
				WHEN EarnCode=1 THEN 0
				ELSE 1
			END AS salaried,
			CASE
				WHEN PRDept IN ('031','032','033','034','051','052','053','054','081','082','083','084') THEN '30'
				WHEN PRDept IN ('014','013','043','044') THEN '01'
				WHEN PRDept IN ('023', '024') THEN '02'
				WHEN PRDept IN ('062', '063', '064') THEN '06'
			END AS prdept
		FROM PREH
		WHERE ActiveYN='Y'
		AND PRGroup=1
		AND PRCo=1
		AND PRDept IN ('031','032','033','034','051','052','053','054','081','082','083','084', '014','013','043','044', '023', '024', '062', '063', '064')
	`
	syncAllCraftTemplates = `
		SELECT PRTC.Template AS template, PRTC.Class AS class, PRCC.Description AS description
		FROM PRTC 
		JOIN PRCC ON PRTC.Class = PRCC.Class 
		AND PRTC.PRCo = PRCC.PRCo 
		AND PRTC.Craft = PRCC.Craft 
		WHERE PRTC.Template IN (
			SELECT CraftTemplate 
			FROM JCJM 
			WHERE JCCo=1 
			AND JobStatus=1 
			AND udFMTS='Y'
		)
	`
	syncAllEquipment = `
		SELECT Equipment AS equipment, Description AS description, Department AS department
		FROM EMEM
		WHERE Status='A'
		AND PRCo=1
	`
)

//get results from viewpoint
//epiloge drop sync tables
//epilogue add sync tables and insert viewpoint results
/* compare sync table to table
if record in both do nothing
if record in sync but not table -> add to table
if record in table but not sync -> mark closed
*/
