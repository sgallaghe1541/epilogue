package db

import (
	"database/sql"
	"fmt"
	"slices"
	"time"
)

type TimeCardHeader struct {
	ID           int            `db:"id"`
	Job          sql.NullString `db:"job"`
	Description  sql.NullString `db:"jobdescription"`
	Date         sql.NullTime   `db:"workdate"`
	WEDate       sql.NullTime   `db:"wedate"`
	CreatedBy    int            `db:"createdby"`
	Status       string         `db:"tcstatus"`
	LastModified time.Time      `db:"lastmodified"`
	ModifiedBy   int            `db:"modifiedby"`
}

type TimeCardEmployee struct {
	ID       int            `db:"tceid"`
	TCHID    int            `db:"tchid"`
	Employee string         `db:"employee"`
	Name     string         `db:"fullname"`
	Date     time.Time      `db:"workdate"`
	Job      string         `db:"job"`
	Phase    string         `db:"phase"`
	Class    sql.NullString `db:"class"`
	PayCode  string         `db:"paycode"`
	Hours    float64        `db:"tcehours"`
}

type TimeCardEmployees []TimeCardEmployee

func (emps TimeCardEmployees) GetUniqueEmployees() ([]string, []string) {
	nums := []string{}
	names := []string{}
	for _, emp := range emps {
		if slices.Contains(nums, emp.Employee) {
			continue
		} else {
			nums = append(nums, emp.Employee)
			nameString := fmt.Sprintf("%s -- %s", emp.Employee, emp.Name)
			names = append(names, nameString)
		}
	}
	return nums, names
}

func (emps TimeCardEmployees) GetPhases() map[int]string {
	phases := []string{}
	for _, emp := range emps {
		if slices.Contains(phases, emp.Phase) {
			continue
		} else {
			phases = append(phases, emp.Phase)
		}
	}
	phaseMap := map[int]string{}
	for i, phase := range phases {
		phaseMap[i] = phase
	}
	return phaseMap
}

func (e *EpilogueConnection) GetTimecardsByUser(userid int, timeCardStatus string) ([]TimeCardHeader, error) {
	args := map[string]interface{}{"userid": userid, "status": timeCardStatus}
	timecards := []TimeCardHeader{}

	err := e.DB.Select(&timecards, `
		SELECT *
		FROM time_card_headers
		WHERE createdby = :userid
		AND tcstatus = :status
		`, args)
	if err != nil {
		return nil, err
	}
	return timecards, nil
}

func (e *EpilogueConnection) GetTimecardIDsByUser(userid int, timeCardStatus string) ([]TimeCardHeader, error) {
	timecards := []TimeCardHeader{}

	err := e.DB.Select(&timecards, `
		SELECT id
		FROM time_card_headers
		WHERE createdby = ?
		AND tcstatus = ?
		`, userid, timeCardStatus)
	if err != nil {
		return nil, err
	}
	return timecards, nil
}

func (e *EpilogueConnection) GetTimecardByID(tcID int) (TimeCardHeader, error) {
	timecard := TimeCardHeader{}

	err := e.DB.Get(&timecard, `
		SELECT *
		FROM time_card_headers
		WHERE id = ?
		`, tcID)
	if err != nil {
		return timecard, err
	}
	return timecard, nil
}

func (e *EpilogueConnection) GetTimecardEmployees(timeCardHeaderID int) ([]TimeCardEmployee, error) {
	employees := []TimeCardEmployee{}

	err := e.DB.Select(&employees, `
		SELECT *
		FROM time_card_employees
		WHERE tchid = ?`, timeCardHeaderID)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (e *EpilogueConnection) LoadTimecard(timeCardHeaderID int) (TimeCardHeader, []TimeCardEmployee, error) {
	header, err := e.GetTimecardByID(timeCardHeaderID)
	if err != nil {
		return header, nil, err
	}
	employees, err := e.GetTimecardEmployees(timeCardHeaderID)
	if err != nil {
		return header, employees, err
	}
	return header, employees, nil
}

func (e *EpilogueConnection) NewTimecard(userid int) (int64, error) {

	vals := map[string]interface{}{
		"createdby":    userid,
		"tcstatus":     "new",
		"lastmodified": time.Now(),
		"modifiedby":   userid,
	}

	result, err := e.DB.NamedExec(`
		INSERT INTO time_card_headers (
			createdby, tcstatus, lastmodified, modifiedby)
		VALUES (
			:createdby, :tcstatus, :lastmodified, :modifiedby
		)`, vals)
	if err != nil {
		return 0, err
	}
	tchid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return tchid, nil
}
