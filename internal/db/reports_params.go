package db

import (
	"database/sql"
)

type EpilogueReport struct {
	ID   int    `db:"reportid"`
	Name string `db:"reportname"`
	URL  string `db:"reporturl"`
}

func (e *EpilogueReport) VPURL() string {
	return "/vp/" + e.URL
}

type ReportParameter struct {
	ID          sql.NullString `db:"paramid"`
	Name        string         `db:"paramname"`
	Description string         `db:"paramdesc"`
	Type        string         `db:"paramtype"`
	ReportID    string         `db:"reportid"`
	SelectList  []*SelectOption
}

type SelectOption struct {
	Value     string `db:""`
	DisplayAs string `db:""`
}
