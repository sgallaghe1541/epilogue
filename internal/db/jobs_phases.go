package db

import (
	"database/sql"
	"fmt"
)

const (
	jobListQuery = `
		SELECT job, description, state, certified, template
		FROM Jobs
		WHERE active = 1
		AND job LIKE ?
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

func (e *EpilogueConnection) GetJobsByDivision(jobEnding string) ([]Job, error) {
	jobs := []Job{}

	err := e.DB.Select(&jobs, jobListQuery, jobEnding)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

const (
	phasesByJob = `
		SELECT job, phase, description
		FROM Phase
		AND job = ?
	`
)

type Phase struct {
	Job         string `db:"job"`
	Phase       string `db:"phase"`
	Description string `db:"description"`
}

func (p Phase) SelectValue() string {
	return fmt.Sprintf("%s -- %s", p.Phase, p.Description)
}

func (p Phase) SelectString() string {
	return fmt.Sprintf("%s -- %s", p.Phase, p.Description)
}

func (e *EpilogueConnection) GetPhasesByJob(job Job) ([]Phase, error) {
	phases := []Phase{}

	err := e.DB.Select(&phases, phasesByJob, job)
	if err != nil {
		return nil, err
	}
	return phases, nil
}
