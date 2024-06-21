package viewpoint

type Division struct {
	JobCostDepts []string
	PayrollDepts []string
}

var Grading = Division{
	JobCostDepts: []string{"1", "4"},
	PayrollDepts: []string{"014", "013", "043", "044"},
}

var Paving = Division{
	JobCostDepts: []string{"2", "20", "3", "33", "5", "55", "8", "88"},
	PayrollDepts: []string{"023", "024", "031", "032", "033", "034", "051", "052", "053", "054", "081", "082", "083", "084", "113", "116", "133"},
}
