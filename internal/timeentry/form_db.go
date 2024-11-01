package timeentry

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/utils"
)

func TimeCardFormFromDB(tc db.TimeCardHeader, epilogue *db.EpilogueConnection) (*TimeCardForm, error) {
	form := NewTimeCardForm(tc.ID)

	if tc.Job.Valid {
		job, err := epilogue.GetJobByNumber(tc.Job.String)
		if err != nil {
			return nil, err
		}
		form.Job = &job
	}
	if tc.Date.Valid {
		form.WorkDate = tc.Date.Time
	}

	phaseMap := tc.GetPhases()
	empNums, empNames := tc.GetUniqueEmployees()

	if len(empNums) != len(empNames) {
		return nil, errors.New("count of employee numbers does not equal count of employee names")
	}

	employeeDetails := map[int]*TimeCardEmployeeForm{}

	for i := 0; i < len(empNums); i++ {
		class, earn, hours := getEmployeeDetails(tc.Employees, empNums[i], phaseMap)
		employeeDetails[i] = &TimeCardEmployeeForm{
			EmpNameString: empNames[i],
			Class:         class,
			Earn:          earn,
			Hours:         hours,
		}
	}

	form.CountPhases = len(phaseMap)
	form.CountEmployees = len(empNums)
	form.Phases = phaseMap
	form.Employees = employeeDetails

	return form, nil
}

func getEmployeeDetails(emps []db.TimeCardEmployee, empNum string, phases map[int]string) (string, string, map[int]float64) {
	class := ""
	earn := ""
	hours := map[int]float64{}

	for _, emp := range emps {
		if emp.Employee != empNum {
			continue
		}
		if earn == "" && emp.PayCode.Valid {
			earn = emp.PayCode.String
		}
		key, err := getPhaseKey(emp.Phase.String, phases)
		if err != nil {
			return "", "", nil
		}
		hours[key] = emp.Hours.Float64
	}

	return class, earn, hours
}

func getPhaseKey(p string, phases map[int]string) (int, error) {
	for k, phase := range phases {
		if phase == p {
			return k, nil
		}
	}
	return 0, fmt.Errorf("failed to find phase: %s in phases", p)
}

func (f *TimeCardForm) ToDB(epilogue *db.EpilogueConnection, status string, userID int) (*db.TimeCardHeader, error) {

	var j sql.NullString

	if f.Job != nil {
		j = sql.NullString{String: f.Job.Job, Valid: true}
	} else {
		j = sql.NullString{Valid: false}
	}

	tcHeader := &db.TimeCardHeader{
		ID:           f.ID,
		Date:         utils.MakeSQLNullTime(f.WorkDate),
		Job:          j,
		Status:       status,
		LastModified: time.Now(),
		ModifiedBy:   userID,
		Employees:    []db.TimeCardEmployee{},
	}

	keys := []int{}
	for k := range f.Employees {
		keys = append(keys, k)
	}

	keys = utils.SortIntKeys(keys)

	for _, k := range keys {
		if emp, ok := f.Employees[k]; ok {
			rows, err := emp.toDBEmployeeRow(f.ID, f.WorkDate, f.Job.Job, f.Phases)
			if err != nil {
				return nil, err
			}
			tcHeader.Employees = append(tcHeader.Employees, rows...)
		}
	}

	return tcHeader, nil
}

func (e *TimeCardEmployeeForm) toDBEmployeeRow(id int, date time.Time, job string, phases map[int]string) ([]db.TimeCardEmployee, error) {
	empRows := []db.TimeCardEmployee{}

	empNum, empName, err := utils.Split(e.EmpNameString)
	if err != nil {
		return empRows, err
	}

	keys := utils.SortIntKeysS(phases)

	for _, k := range keys {
		for num, hours := range e.Hours {
			if k == num {
				row := db.TimeCardEmployee{
					TCHID:    id,
					Employee: empNum,
					Name:     empName,
					Date:     utils.MakeSQLNullTime(date),
					Job:      utils.MakeSQLNullString(job),
					Phase:    utils.MakeSQLNullString(phases[k]),
					Class:    utils.MakeSQLNullString(e.Class),
					PayCode:  utils.MakeSQLNullString(e.Earn),
					Hours:    utils.MakeSQLNullFloat64(hours),
				}
				empRows = append(empRows, row)
			}
		}
	}

	return empRows, nil

}
