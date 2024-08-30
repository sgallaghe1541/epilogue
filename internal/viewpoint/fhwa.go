package viewpoint

import "fmt"

const (
	FHWAjobQuery string = `
		SELECT
			s.classDescription AS classDescription,
			s.raceDescription AS raceDescription,
			s.sex AS sex,
			COUNT(s.employee) AS numEmployees
		FROM (
			SELECT
				PRTH.Employee    AS employee,
				PREH.Class       AS class,
				PRCC.Description AS classDescription,
				PREH.Race        AS race,
				PRRC.Description AS raceDescription,
				PREH.Sex         AS sex
			FROM PREH
			JOIN PRTH
			ON PRTH.PRCo = PREH.PRCo
			AND PRTH.Employee = PREH.Employee
			JOIN PRCC 
			ON PREH.Craft = PRCC.Craft AND PREH.Class = PRCC.Class AND PREH.PRCo = PRCC.PRCo
			LEFT JOIN JCJM 
			ON PRTH.Job = JCJM.Job
			JOIN PRRC
			ON PREH.PRCo = PRRC.PRCo AND PREH.Race = PRRC.Race
			WHERE PRTH.PRCo = 1
			AND PRTH.PREndDate = :wedate
			AND PRTH.Hours <> 0
			AND JCJM.Job = :job
			GROUP BY 
				PRTH.Employee,
				PREH.Class,
				PRCC.Description,
				PREH.Race,
				PRRC.Description,
				PREH.Sex
			) s
		GROUP BY 
			s.classDescription,
			s.raceDescription,
			s.sex
	`

	FHWAgroupQuery string = `
		SELECT
			s.classDescription AS classDescription,
			s.raceDescription AS raceDescription,
			s.sex AS sex,
			COUNT(s.employee) AS numEmployees
		FROM (
			SELECT
				PRTH.Employee    AS employee,
				PREH.Class       AS class,
				PRCC.Description AS classDescription,
				PREH.Race        AS race,
				PRRC.Description AS raceDescription,
				PREH.Sex         AS sex
			FROM PREH
			JOIN PRTH
			ON PRTH.PRCo = PREH.PRCo
			AND PRTH.Employee = PREH.Employee
			JOIN PRCC 
			ON PREH.Craft = PRCC.Craft AND PREH.Class = PRCC.Class AND PREH.PRCo = PRCC.PRCo
			JOIN PRRC
			ON PREH.PRCo = PRRC.PRCo AND PREH.Race = PRRC.Race
			WHERE PRTH.PRCo = 1
			AND PRTH.PREndDate = :wedate
			AND PRTH.Hours <> 0
			AND PREH.PRGroup IN (:prgroups)
			GROUP BY 
				PRTH.Employee,
				PREH.Class,
				PRCC.Description,
				PREH.Race,
				PRRC.Description,
				PREH.Sex
			) s
		GROUP BY 
			s.classDescription,
			s.raceDescription,
			s.sex
	`
)

type FHWAClassification struct {
	Class string `db:"classDescription"`
	Race  string `db:"raceDescription"`
	Sex   string `db:"sex"`
	Count int    `db:"numEmployees"`
}

func (f FHWAClassification) raceGender() string {
	var gender string

	if f.Sex == "M" {
		gender = "Male"
	} else if f.Sex == "F" {
		gender = "Female"
	} else {
		gender = ""
	}

	return fmt.Sprintf("%s %s", f.Race, gender)
}

type FHWAClassificationResult struct {
	raceSex map[string][]string
	headers []string
	Result  []*FHWAClassification
}

func (f *FHWAClassificationResult) setHeaders() error {
	if f.raceSex != nil {
		return fmt.Errorf("fhwa headers have already been set")
	}
	if f.Result == nil {
		return fmt.Errorf("no fhwa results found")
	}

	headers := make(map[string][]string)

	for _, result := range f.Result {
		if headers[result.Race] == nil {
			headers[result.Race] = []string{"Male", "Female"}
		}
	}
	f.raceSex = headers

	return nil
}

func (f *FHWAClassificationResult) Headers() []string {
	if f.raceSex == nil {
		err := f.setHeaders()
		if err != nil {
			return nil
		}
	}
	headers := []string{"Classification"}

	for race, sexes := range f.raceSex {
		for _, sex := range sexes {
			headers = append(headers, fmt.Sprintf("%s %s", race, sex))
		}
	}
	f.headers = headers
	return headers
}

func (f *FHWAClassificationResult) Classes() []string {
	classes := []string{}

	for _, r := range f.Result {
		var found bool = false
		for _, class := range classes {
			if r.Class == class {
				found = true
			}
		}
		if !found {
			classes = append(classes, r.Class)
		}
	}
	return classes
}

func (f *FHWAClassificationResult) Data() [][]string {
	headers := f.headers
	classes := f.Classes()
	data := make([][]string, len(classes))

	for i, class := range classes {
		classData := make([]string, len(headers))
		classData[0] = class
		for _, r := range f.Result {
			if r.Class == class {
				for j := 1; j < len(headers); j++ {
					if r.raceGender() == headers[j] {
						classData[j] = fmt.Sprintf("%d", r.Count)
					}
				}
			}
		}
		data[i] = classData
	}

	return data
}
