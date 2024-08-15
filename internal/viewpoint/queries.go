package viewpoint

const (
	JobListQuery = `
		SELECT Job AS job, Description AS description
		FROM JCJM 
		WHERE JCCo=1 
		AND Job LIKE :jobending 
		AND JobStatus=1 
		AND udFMTS='Y'
	`
)

type Excelable interface {
	Headers() *[]interface{}
	ToExcel(string) error
}
