package timeentry

import (
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
	CountPhases    string
	CountEmployees string
	CountEquipment string
	Phases         []string
	Employees      [][]string
	Equipment      [][]string
	FieldErrors    map[string]string
}
