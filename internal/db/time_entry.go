package db

import (
	"database/sql"
	"time"
)

type TimeCardHeader struct {
	ID           int       `db:"id"`
	Job          string    `db:"job"`
	Description  string    `db:"jobdescription"`
	Date         time.Time `db:"workdate"`
	WEDate       time.Time `db:"wedate"`
	CreatedBy    int       `db:"createdby"`
	Status       string    `db:"tcstatus"`
	LastModified time.Time `db:"lastmodified"`
	ModifiedBy   int       `db:"modifiedby"`
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
