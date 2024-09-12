package timeentry

import (
	"time"
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
