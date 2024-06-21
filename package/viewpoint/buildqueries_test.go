package viewpoint

import (
	"testing"
)

func TestBuildQueries(t *testing.T) {
	expectedQuery := `
		SELECT JCJM.Job AS Job, JCJM.Description AS Description
		FROM JCJM 
		JOIN JCCM ON JCJM.JCCo = JCCM.JCCo AND JCJM.Contract = JCCM.Contract
		WHERE JCJM.JCCo=1
		AND JCCM.Department IN(?, ?, ?, ?, ?, ?, ?, ?)
		AND JCJM.JobStatus=1
	`
	// expectedArgs := []interface{}{"1"," 4"}

	query, _, err := BuildInQuery(JobList, Paving)
	if err != nil {
		t.Fatalf("BuildQuery didn't work at all: %s", err)
	}

	if expectedQuery != query {
		t.Fatalf("Wanted %s. \nGot %s", expectedQuery, query)
	}
}
