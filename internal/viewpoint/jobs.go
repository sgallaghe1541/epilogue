package viewpoint

import (
	"database/sql"
	"fmt"
)

const (
	jobListQuery = `
		SELECT Job AS job, Description AS description, PRStateCode AS state, Certified AS certified, CraftTemplate AS template
		FROM JCJM 
		WHERE JCCo=1 
		AND Job LIKE ?
		AND JobStatus=1 
		AND udFMTS='Y'
	`
	phasesByJob = `
		SELECT Phase AS phase, Description AS description 
		FROM JCJP
		WHERE JCCo = 1
		AND udFMTS = 'Y'
		AND Job = ?
	`
)

type Job struct {
	Job           string         `db:"job"`
	Description   sql.NullString `db:"description"`
	State         string         `db:"state"`
	Certified     string         `db:"certified"`
	CraftTemplate sql.NullInt64  `db:"template"`
}

func (j Job) SelectValue() string {
	return fmt.Sprintf("%s -- %s", j.Job, j.Description.String)
}

func (j Job) SelectString() string {
	return fmt.Sprintf("%s -- %s", j.Job, j.Description.String)
}

type Phase struct {
	Phase       string `db:"phase"`
	Description string `db:"description"`
}

func (p Phase) SelectValue() string {
	return fmt.Sprintf("%s -- %s", p.Phase, p.Description)
}

func (p Phase) SelectString() string {
	return fmt.Sprintf("%s -- %s", p.Phase, p.Description)
}

func (v *ViewpointConnection) GetJobsByDivision(jobEnding string) ([]Job, error) {
	jobs := []Job{}

	err := v.DB.Select(&jobs, jobListQuery, jobEnding)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (v *ViewpointConnection) GetPhasesByJob(job string) ([]Phase, error) {
	phases := []Phase{}

	err := v.DB.Select(&phases, phasesByJob, job)
	if err != nil {
		return nil, err
	}
	return phases, nil
}
