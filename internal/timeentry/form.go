package timeentry

import (
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
)

type TimeCardForm struct {
	ID             int
	WorkDate       time.Time
	Job            *db.Job
	CountPhases    int
	CountEmployees int
	CountEquipment int
	Phases         map[int]string
	Employees      map[int]*TimeCardEmployeeForm
	Equipment      [][]string
	FieldErrors    map[string]string
}

type TimeCardEmployeeForm struct {
	EmpNameString string
	Class         string
	Earn          string
	Hours         map[int]float64
}

func (emp TimeCardEmployeeForm) Selected() bool {
	return emp.EmpNameString != ""
}

func NewTimeCardForm(id int) *TimeCardForm {
	return &TimeCardForm{
		ID:          id,
		Phases:      map[int]string{},
		Employees:   map[int]*TimeCardEmployeeForm{},
		FieldErrors: map[string]string{},
	}
}
