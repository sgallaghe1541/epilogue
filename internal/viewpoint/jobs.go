package viewpoint

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

const (
	jobListQuery = `
		SELECT Job AS job, Description AS description
		FROM JCJM 
		WHERE JCCo=1 
		AND Job LIKE ?
		AND JobStatus=1 
		AND udFMTS='Y'
	`
)

type Job struct {
	Job         string         `db:"job"`
	Description sql.NullString `db:"description"`
}

type JobModel struct {
	DB *sqlx.DB
}

func (j *JobModel) GetJobsByDivision(jobEnding string) ([]Job, error) {
	jobs := []Job{}

	err := j.DB.Select(&jobs, jobListQuery, jobEnding)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
