package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type TimecardModel struct {
	DB *sqlx.DB
}

type TimecardHeader struct {
	ID           int       `db:"id"`
	Job          string    `db:"job"`
	Date         time.Time `db:"workdate"`
	WEDate       time.Time `db:"wedate"`
	CreatedBy    int       `db:"createdby"`
	Status       string    `db:"tcstatus"`
	LastModified time.Time `db:"lastmodified"`
	ModifiedBy   int       `db:"modifiedby"`
}

type TimecardEmployee struct {
	ID      int            `db:"tceid"`
	TCHID   int            `db:"tchid"`
	Date    time.Time      `db:"workdate"`
	Class   sql.NullString `db:"class"`
	Name    string         `db:"fullname"`
	PayCode string         `db:"paycode"`
	Phase   string         `db:"phase"`
	Hours   float64        `db:"tcehours"`
}

func (t *TimecardModel) GetTimecards(userid int, timeCardStatus string) ([]TimecardHeader, error) {
	args := map[string]interface{}{"userid": userid, "status": timeCardStatus}
	timecards := []TimecardHeader{}

	err := t.DB.Select(&timecards, `
		SELECT * FROM time_card_header
		WHERE createdby = :userid
		AND tcstatus = :status
		`, args)
	if err != nil {
		return nil, err
	}
	return timecards, nil
}

func (t *TimecardModel) NewTimecard(job string, date time.Time, userid int) (int64, error) {
	wedate := getWEDate(date)

	vals := map[string]interface{}{
		"job":          job,
		"workdate":     date,
		"wedate":       wedate,
		"createdby":    userid,
		"tcstatus":     "draft",
		"lastmodified": time.Now(),
		"modifiedby":   userid,
	}

	result, err := t.DB.NamedExec(`
		INSERT INTO time_card_header (
			job, workdate, wedate, createdby, 
			tcstatus, lastmodified, modifiedby)
		VALUES (
			:job, :workdate, :wedate, :createdby, 
			:tcstatus, :lastmodified, :modifiedby
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

func (t *TimecardModel) GetTimecardEmployees(timeCardHeaderID int) ([]TimecardEmployee, error) {
	args := map[string]interface{}{"tchid": timeCardHeaderID}
	employees := []TimecardEmployee{}

	err := t.DB.Select(&employees, `
		SELECT * FROM time_card_employees
		WHERE tchid = :tchid`, args)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func getWEDate(date time.Time) time.Time {
	for date.Weekday() != 6 {
		date = date.Add(time.Hour * 24)
	}
	return date
}
