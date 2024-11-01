package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
)

type TimeCardHeader struct {
	ID           int            `db:"id"`
	Date         sql.NullTime   `db:"workdate"`
	Job          sql.NullString `db:"job"`
	CreatedBy    int            `db:"createdby"`
	Status       string         `db:"tcstatus"`
	LastModified time.Time      `db:"lastmodified"`
	ModifiedBy   int            `db:"modifiedby"`
	Employees    []TimeCardEmployee
}

type TimeCardEmployee struct {
	ID       int             `db:"tceid"`
	TCHID    int             `db:"tchid"`
	Employee string          `db:"employee"`
	Name     string          `db:"fullname"`
	Date     sql.NullTime    `db:"workdate"`
	Job      sql.NullString  `db:"job"`
	Phase    sql.NullString  `db:"phase"`
	Class    sql.NullString  `db:"class"`
	PayCode  sql.NullString  `db:"paycode"`
	Hours    sql.NullFloat64 `db:"tcehours"`
}

func (t *TimeCardHeader) GetEmployeesFromDB(data *sqlx.DB, logger *slog.Logger) error {
	query := `
		SELECT * 
		FROM time_card_employees
		WHERE tchid = ?
	`
	err := data.Select(&t.Employees, query, fmt.Sprintf("%d", t.ID))
	if err != nil {
		logger.Error("failed to load employees", "tchid", t.ID, "error", err.Error())
		return err
	}
	return nil
}

func (tc TimeCardHeader) GetUniqueEmployees() ([]string, []string) {
	nums := []string{}
	names := []string{}
	for _, emp := range tc.Employees {
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

func (tc TimeCardHeader) GetPhases() map[int]string {
	phases := []string{}
	for _, emp := range tc.Employees {
		if emp.Phase.Valid {
			if slices.Contains(phases, string(emp.Phase.String)) {
				continue
			} else {
				phases = append(phases, string(emp.Phase.String))
			}
		}
	}
	phaseMap := map[int]string{}
	for i, phase := range phases {
		phaseMap[i] = phase
	}
	return phaseMap
}

func (tc TimeCardHeader) UpdateTimeCard(epilogue *EpilogueConnection, logger *slog.Logger) error {
	updateHeaderStmt := `
		UPDATE time_card_headers
		SET 
			job = :job,
			workdate = :workdate,
			lastmodified = :lastmodified,
			modifiedby = :modifiedby,
			tcstatus = 'draft'
		WHERE id = :id
	`
	deleteEmployeeStmt := `DELETE FROM time_card_employees WHERE tchid = ?`
	insertEmployeeStmt := `
		INSERT INTO time_card_employees (tchid, employee, fullname, workdate, job, phase, class, paycode, tcehours)
		VALUES (:tchid, :employee, :fullname, :workdate, :job, :phase, :class, :paycode, :tcehours)
	`

	tx, err := epilogue.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(updateHeaderStmt, tc)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			logger.Error("failed to rollback update timecard employees after delete error", "err", rollErr.Error())
		}
		return err
	}
	_, err = tx.Exec(deleteEmployeeStmt, fmt.Sprintf("%d", tc.ID))
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			logger.Error("failed to rollback update timecard employees after delete error", "err", rollErr.Error())
		}
		return err
	}

	_, err = tx.NamedExec(insertEmployeeStmt, tc.Employees)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			logger.Error("failed to rollback update timecard employees after insert error", "err", rollErr.Error())
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (e *EpilogueConnection) GetTimecardsByUserStatus(userid, timeCardStatus string) ([]TimeCardHeader, error) {
	timecards := []TimeCardHeader{}

	err := e.DB.Select(&timecards, `
		SELECT *
		FROM time_card_headers
		WHERE createdby = ?
		AND tcstatus = ?
		`, userid, timeCardStatus)
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

func (e *EpilogueConnection) GetTimecardByID(tcID string) (TimeCardHeader, error) {
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

func (e *EpilogueConnection) UpdateTimecardHeader(timecard TimeCardHeader) error {
	_, err := e.DB.NamedExec(`
		UPDATE time_card_headers
		SET 
			job = :job,
			workdate = :workdate,
			wedate = :wedate,
			tcstatus = :tcstatus,
			lastmodified = :lastmodified
		WHERE id = :id
	`, timecard)

	if err != nil {
		return err
	}
	return nil
}

func (e *EpilogueConnection) LoadTimecard(timeCardHeaderID string, logger *slog.Logger) (TimeCardHeader, error) {
	header, err := e.GetTimecardByID(timeCardHeaderID)
	if err != nil {
		logger.Error("failed to load time card", "id", timeCardHeaderID)
		return header, err
	}
	err = header.GetEmployeesFromDB(e.DB, logger)
	if err != nil {
		return header, err
	}
	return header, nil
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
