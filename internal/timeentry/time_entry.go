package timeentry

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type TimeCardHeaderCreateForm struct {
	WorkDate    time.Time
	Job         string
	FieldErrors map[string]string
}

type TimeCardEmployeesForm struct {
	Phases      string
	FieldErrors map[string]string
}

type TimeCardModel struct {
	DB *sqlx.DB
}

type TimeCardHeader struct {
	ID           int       `db:"id"`
	Job          string    `db:"job"`
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

func (t *TimeCardModel) GetTimecardsByUser(userid int, timeCardStatus string) ([]TimeCardHeader, error) {
	args := map[string]interface{}{"userid": userid, "status": timeCardStatus}
	timecards := []TimeCardHeader{}

	err := t.DB.Select(&timecards, `
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

func (t *TimeCardModel) GetTimecardByID(tcID int) (TimeCardHeader, error) {
	timecard := TimeCardHeader{}

	err := t.DB.Get(&timecard, `
		SELECT *
		FROM time_card_headers
		WHERE id = ?
		`, tcID)
	if err != nil {
		return timecard, err
	}
	return timecard, nil
}

func (t *TimeCardModel) NewTimecard(job string, date time.Time, userid int) (int64, error) {
	wedate := getWEDate(date)

	vals := map[string]interface{}{
		"job":          job,
		"workdate":     date,
		"wedate":       wedate,
		"createdby":    userid,
		"tcstatus":     "new",
		"lastmodified": time.Now(),
		"modifiedby":   userid,
	}

	result, err := t.DB.NamedExec(`
		INSERT INTO time_card_headers (
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

func (t *TimeCardModel) GetTimecardEmployees(timeCardHeaderID int) ([]TimeCardEmployee, error) {
	employees := []TimeCardEmployee{}

	err := t.DB.Select(&employees, `
		SELECT *
		FROM time_card_employees
		WHERE tchid = ?`, timeCardHeaderID)
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
