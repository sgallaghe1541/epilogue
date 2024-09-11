package viewpoint

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
	payPeriodStatus = `
		SELECT Status AS status
		FROM PRPC
		WHERE PRCo = 1
		AND PREndDate = ?
	`
	employeesByDivision = `
		SELECT Employee AS employee, CONCAT(FirstName, " ", LastName) AS name
		FROM PREH
		WHERE ActiveYN = 'Y'
		AND PRGroup = 1
		AND PRCo = 1
		AND PRDept IN (:payrolldepts)
	`
	phasesByJob = `
		SELECT Phase AS phase, Description AS description 
		FROM JCJP
		WHERE JCCo = 1
		AND udFMTS = 'Y'
		AND Job = ?
	`
)

type SelectOption interface {
	SelectValue() string
	SelectString() string
}

type PayPeriod struct {
	Status int `db:"status"`
}

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

type Employee struct {
	Employee string `db:"employee"`
	Name     string `db:"name"`
}

func (e Employee) SelectValue() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
}

func (e Employee) SelectString() string {
	return fmt.Sprintf("%s -- %s", e.Employee, e.Name)
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

type TimeEntryModel struct {
	DB *sqlx.DB
}

func (t *TimeEntryModel) GetJobsByDivision(jobEnding string) ([]Job, error) {
	jobs := []Job{}

	err := t.DB.Select(&jobs, jobListQuery, jobEnding)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (t *TimeEntryModel) GetPayPeriodStatus(date time.Time) int {
	pp := PayPeriod{}

	err := t.DB.Get(&pp, payPeriodStatus, date.Format("01/02/06"))
	if err != nil {
		return 0
	}
	return pp.Status
}

func (t *TimeEntryModel) GetEmployeesByDivision(div QueryArgs) ([]Employee, error) {
	emps := []Employee{}

	query, args, err := BuildInQuery(employeesByDivision, div)
	if err != nil {
		return nil, err
	}
	err = t.DB.Select(&emps, query, args)
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func (t *TimeEntryModel) GetPhasesByJob(job string) ([]Phase, error) {
	phases := []Phase{}

	err := t.DB.Select(&phases, phasesByJob, job)
	if err != nil {
		return nil, err
	}
	return phases, nil
}
