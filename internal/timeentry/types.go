package timeentry

import (
	"errors"
	"fmt"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
)

type SelectOption interface {
	SelectValue() string
	SelectString() string
}

type TimeCardHeaderForm struct {
	WorkDate    time.Time
	Job         *db.Job
	FieldErrors map[string]string
}

type TimeCardDetailForm struct {
	CountPhases    int
	CountEmployees int
	CountEquipment int
	Phases         map[int]string
	Employees      map[int]TimeCardEmployeeForm
	Equipment      [][]string
	FieldErrors    map[string]string
}

type TimeCardEmployeeForm struct {
	EmpNameString string
	Class         string
	Earn          string
	Hours         map[string]string
}

func (emp TimeCardEmployeeForm) Selected() bool {
	return emp.EmpNameString != ""
}

func PopulateTimeCardHeaderForm(tc db.TimeCardHeader, epilogue *db.EpilogueConnection) (*TimeCardHeaderForm, error) {
	if !tc.Job.Valid {
		if !tc.Date.Valid {
			form := TimeCardHeaderForm{}
			return &form, nil
		} else {
			form := TimeCardHeaderForm{
				WorkDate: tc.Date.Time,
			}
			return &form, nil
		}
	} else {
		job, err := epilogue.GetJobByNumber(tc.Job.String)
		if err != nil {
			return nil, err
		}
		if !tc.Date.Valid {
			form := TimeCardHeaderForm{
				Job: &job,
			}
			return &form, nil
		} else {
			form := TimeCardHeaderForm{
				Job:      &job,
				WorkDate: tc.Date.Time,
			}
			return &form, nil
		}
	}
}

func PopulateTimeCardDetailForm(emps db.TimeCardEmployees) (*TimeCardDetailForm, error) {
	phaseMap := emps.GetPhases()
	empNums, empNames := emps.GetUniqueEmployees()

	if len(empNums) != len(empNames) {
		return nil, errors.New("count of employee numbers does not equal count of employee names")
	}

	employeeDetails := map[int]TimeCardEmployeeForm{}

	for i := 0; i < len(empNums); i++ {
		class, earn, hours := getEmployeeDetails(emps, empNums[i])
		employeeDetails[i] = TimeCardEmployeeForm{
			EmpNameString: empNames[i],
			Class:         class,
			Earn:          earn,
			Hours:         hours,
		}
	}

	form := TimeCardDetailForm{
		CountPhases:    len(phaseMap),
		CountEmployees: len(empNums),
		Phases:         phaseMap,
		Employees:      employeeDetails,
	}
	return &form, nil
}

func getEmployeeDetails(emps db.TimeCardEmployees, empNum string) (string, string, map[string]string) {
	class := ""
	earn := ""
	hours := map[string]string{}

	for _, emp := range emps {
		if emp.Employee != empNum {
			continue
		}
		if earn == "" {
			earn = emp.PayCode
		}
		hours[emp.Phase] = fmt.Sprintf("%.2f", emp.Hours)
	}

	return class, earn, hours
}
