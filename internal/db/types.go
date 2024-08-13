package db

type EpilogueReport struct {
	ID         int    `db:"reportid"`
	Name       string `db:"reportname"`
	URL        string `db:"reporturl"`
	Permission int    `db:"permission"`
}

func (e *EpilogueReport) VPURL() string {
	return "/vp/" + e.URL
}

type ReportParameter struct {
	Name        string `db:"paramname"`
	Description string `db:"paramdesc"`
	Type        string `db:"paramtype"`
	ReportID    string `db:"reportid"`
}

type Division struct {
	ID          string `db:"divisionid"`
	Description string `db:"description"`
}
