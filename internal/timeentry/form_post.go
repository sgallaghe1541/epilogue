package timeentry

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/utils"
)

func TimeCardFormFromPostForm(id int, postForm map[string][]string, epilogue *db.EpilogueConnection) (*TimeCardForm, error) {

	form := NewTimeCardForm(id)

	for key, value := range postForm {
		if key == "job" {
			err := form.parseJob(value, epilogue)
			if err != nil {
				return nil, err
			}

		} else if key == "workdate" {
			err := form.parseWorkDate(value)
			if err != nil {
				return nil, err
			}

		} else if key == "phasecount" {
			err := form.parsePhaseCount(value)
			if err != nil {
				return nil, err
			}

		} else if strings.HasPrefix(key, "phase-") {
			err := form.parsePhase(key, value)
			if err != nil {
				return nil, err
			}

		} else if key == "employeecount" {
			err := form.parseEmployeeCount(value)
			if err != nil {
				return nil, err
			}

		} else if strings.HasPrefix(key, "employee-") {
			err := form.parseEmployees(key, value)
			if err != nil {
				return nil, err
			}
		} else {
			continue
		}
	}

	return form, nil
}

func hasValue(value []string) bool {
	if len(value) != 1 {
		return false
	}
	if value[0] == "" {
		return false
	}
	return true
}

func getStringValue(value []string) string {
	if hasValue(value) {
		return value[0]
	}
	return ""
}

func getValidFloatValue(value string) (float64, error) {
	val, err := strconv.ParseFloat(value, 64)
	return val, err
}

func (f *TimeCardForm) parseJob(value []string, epilogue *db.EpilogueConnection) error {
	if hasValue(value) {
		num, _, err := utils.Split(value[0])
		if err != nil {
			return err
		}

		job, err := epilogue.GetJobByNumber(num)

		if err != nil {
			return err
		}
		f.Job = &job
	}
	return nil
}

func (f *TimeCardForm) parseWorkDate(value []string) error {
	if hasValue(value) {
		date, err := time.Parse("2006-01-02", value[0])
		if err != nil {
			return err
		}
		f.WorkDate = date
	}
	return nil
}

func (f *TimeCardForm) parsePhaseCount(value []string) error {
	if hasValue(value) {
		count, err := strconv.Atoi(value[0])
		if err != nil {
			return fmt.Errorf("invalid phase count: %s error: %s", value[0], err)
		}
		f.CountPhases = count
	}
	return nil
}

func (f *TimeCardForm) parsePhase(key string, value []string) error {
	if !hasValue(value) {
		return nil
	}

	phaseSlice := strings.Split(key, "-")
	if len(phaseSlice) != 2 {
		return fmt.Errorf("failed to parse phase key: %s", key)
	}
	num, err := strconv.Atoi(phaseSlice[1])
	if err != nil {
		return fmt.Errorf("invalid number in phase key: %s error: %s", key, err)
	}

	phaseString := getStringValue(value)
	if phaseString == "" {
		return nil
	}

	phaseNum, _, err := utils.Split(phaseString)
	if err != nil {
		return fmt.Errorf("could not parse phase number phase: %s error: %s", phaseString, err)
	}

	f.Phases[num] = phaseNum
	return nil
}

func (f *TimeCardForm) parseEmployeeCount(value []string) error {
	if hasValue(value) {
		count, err := strconv.Atoi(value[0])
		if err != nil {
			return fmt.Errorf("invalid employee count: %s error: %s", value[0], err)
		}
		f.CountEmployees = count
	}
	return nil
}

func (f *TimeCardForm) parseEmployees(key string, value []string) error {
	if hasValue(value) {
		employeeSlice := strings.Split(key, "-")
		if len(employeeSlice) == 2 {
			num, err := strconv.Atoi(employeeSlice[1])
			if err != nil {
				return fmt.Errorf("invalid number in employee key: %s error: %s", key, err)
			}
			if tcEmployee, ok := f.Employees[num]; ok {
				tcEmployee.EmpNameString = getStringValue(value)
			} else {
				tcEmployee := &TimeCardEmployeeForm{
					EmpNameString: getStringValue(value),
					Hours:         map[int]float64{},
				}
				f.Employees[num] = tcEmployee
			}

		} else if len(employeeSlice) == 3 {
			if employeeSlice[2] == "class" {
				err := f.parseEmployeeClass(employeeSlice, value)
				if err != nil {
					return err
				}
			} else if employeeSlice[2] == "earn" {
				err := f.parseEmployeeEarn(employeeSlice, value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unknown employee key: %s", key)
			}

		} else if len(employeeSlice) == 4 {
			err := f.parseEmployeeHours(employeeSlice, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *TimeCardForm) parseEmployeeClass(key []string, value []string) error {
	num, err := strconv.Atoi(key[1])
	if err != nil {
		return fmt.Errorf("invalid number in employee key: %s error: %s", key, err)
	}
	if tcEmployee, ok := f.Employees[num]; ok {
		tcEmployee.Class = getStringValue(value)
	} else {
		tcEmployee := &TimeCardEmployeeForm{
			Class: getStringValue(value),
			Hours: map[int]float64{},
		}
		f.Employees[num] = tcEmployee
	}
	return nil
}

func (f *TimeCardForm) parseEmployeeEarn(key, value []string) error {
	num, err := strconv.Atoi(key[1])
	if err != nil {
		return fmt.Errorf("invalid number in employee key: %s error: %s", key, err)
	}
	if tcEmployee, ok := f.Employees[num]; ok {
		tcEmployee.Earn = getStringValue(value)
	} else {
		tcEmployee := &TimeCardEmployeeForm{
			Earn:  getStringValue(value),
			Hours: map[int]float64{},
		}
		f.Employees[num] = tcEmployee
	}
	return nil
}

func (f *TimeCardForm) parseEmployeeHours(key, value []string) error {
	empNum, err := strconv.Atoi(key[1])
	if err != nil {
		return fmt.Errorf("invalid number in employee key: %s error: %s", key, err)
	}

	hrNum, err := strconv.Atoi(key[3])
	if err != nil {
		return fmt.Errorf("invalid hour number in employee key: %s error: %s", key, err)
	}

	if !hasValue(value) {
		//add logging
		return nil
	}

	hours, err := getValidFloatValue(value[0])
	if err != nil {
		return fmt.Errorf("failed to parse hours key: %s, value: %s, error: %s", key, value[0], err)
	}

	if tcEmployee, ok := f.Employees[empNum]; ok {
		tcEmployee.Hours[hrNum] = hours
	} else {
		tcEmployee := &TimeCardEmployeeForm{Hours: map[int]float64{}}
		f.Employees[empNum] = tcEmployee
		tcEmployee.Hours[hrNum] = hours
	}

	return nil
}
