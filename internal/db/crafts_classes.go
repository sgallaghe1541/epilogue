package db

import "fmt"

const (
	classesByJob = `
		SELECT template, class, description
		FROM craft_classes
		WHERE template = ?
	`
)

type CraftClass struct {
	Template    int64  `db:"template"`
	Class       string `db:"class"`
	Description string `db:"description"`
}

func (c CraftClass) SelectValue() string {
	return c.Class
}

func (c CraftClass) SelectString() string {
	return fmt.Sprintf("%s -- %s", c.Class, c.Description)
}

func (e *EpilogueConnection) GetCraftClassByJob(job Job) ([]CraftClass, error) {
	classes := []CraftClass{}

	if !job.CraftTemplate.Valid {
		return classes, fmt.Errorf("%s is not a certified job", job.Job)
	}

	err := e.DB.Select(&classes, classesByJob, job.CraftTemplate.Int64)
	if err != nil {
		return classes, err
	}
	return classes, nil
}
