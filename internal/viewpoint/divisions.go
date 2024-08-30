package viewpoint

type QueryArgs struct {
	JobCostDepts []string `db:"jobcostdepts"`
	PayrollDepts []string `db:"payrolldepts"`
	JobEnding    string   `db:"jobending"`
	StartWEDate  string   `db:"startwedate"`
	EndWEDate    string   `db:"endwedate"`
	Job          string   `db:"job"`
	Employee     string   `db:"employee"`
	Date         string   `db:"wedate"`
	Groups       []string `db:"prgroups"`
}

var Grading = QueryArgs{
	JobCostDepts: []string{"1", "4"},
	PayrollDepts: []string{"014", "013", "043", "044"},
	JobEnding:    "%.01",
}

var Paving = QueryArgs{
	JobCostDepts: []string{"2", "20", "3", "33", "5", "55", "8", "88"},
	PayrollDepts: []string{"023", "024"},
	JobEnding:    "%.02",
}

var Bridge = QueryArgs{
	JobCostDepts: []string{"6"},
	PayrollDepts: []string{"062", "063", "064"},
	JobEnding:    "%.06",
}

var Plants = QueryArgs{
	JobCostDepts: []string{"2", "20", "3", "33", "5", "55", "8", "88"},
	PayrollDepts: []string{"031", "032", "033", "034", "051", "052", "053", "054", "081", "082", "083", "084"},
}

var ChattShop = QueryArgs{
	JobCostDepts: []string{"2", "20", "3", "33", "5", "55", "8", "88"},
	PayrollDepts: []string{"113", "116", "133"},
}
