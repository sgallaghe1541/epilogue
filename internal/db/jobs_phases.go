package db

import (
	"database/sql"
	"fmt"
)

type Job struct {
	Job           string         `db:"job"`
	Description   sql.NullString `db:"description"`
	State         string         `db:"state"`
	Certified     string         `db:"certified"`
	CraftTemplate sql.NullInt64  `db:"template"`
	Department    sql.NullString `db:"department"`
}

func (j Job) SelectValue() string {
	return j.Job
}

func (j Job) SelectString() string {
	return fmt.Sprintf("%s -- %s", j.Job, j.Description.String)
}

func (e *EpilogueConnection) GetJobsByDepartment(jobEnding string) ([]Job, error) {
	query := `
		SELECT job, description, state, certified, template, department
		FROM Jobs
		WHERE active = 1
		AND department = ?
	`
	jobs := []Job{}

	err := e.DB.Select(&jobs, query, jobEnding)
	if err != nil {
		return []Job{}, err
	}
	return jobs, nil
}

func (e *EpilogueConnection) GetJobByNumber(jobNumber string) (Job, error) {
	query := `
		SELECT job, description, state, certified, template, department
		FROM Jobs
		WHERE active = 1
		AND job = ?
	`
	job := Job{}

	err := e.DB.Get(&job, query, jobNumber)
	if err != nil {
		return job, err
	}
	return job, nil
}

type Phase struct {
	Job         string `db:"job"`
	Phase       string `db:"phase"`
	Description string `db:"description"`
}

func (p Phase) SelectValue() string {
	return p.Phase
}

func (p Phase) SelectString() string {
	return fmt.Sprintf("%s -- %s", p.Phase, p.Description)
}

func (e *EpilogueConnection) GetPhasesByJob(job Job) ([]Phase, error) {
	query := `
		SELECT job, phase, description
		FROM Phase
		WHERE job = ?
	`
	phases := []Phase{}

	err := e.DB.Select(&phases, query, job)
	if err != nil {
		return phases, err
	}
	return phases, nil
}
