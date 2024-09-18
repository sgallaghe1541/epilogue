package db

import "fmt"

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
