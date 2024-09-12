package viewpoint

import "time"

const (
	payPeriodStatus = `
		SELECT Status AS status
		FROM PRPC
		WHERE PRCo = 1
		AND PREndDate = ?
	`
)

type PayPeriod struct {
	Status int `db:"status"`
}

func (v *ViewpointConnection) GetPayPeriodStatus(date time.Time) int {
	pp := PayPeriod{}

	err := v.DB.Get(&pp, payPeriodStatus, date.Format("01/02/06"))
	if err != nil {
		return 0
	}
	return pp.Status
}
