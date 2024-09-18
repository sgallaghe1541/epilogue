package sync

import (
	"os"
	"strings"

	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
)

func (s *Syncer) sync() int {
	// get vp templates, jobs, phases, employees, equipment
	// phases := []viewpoint.Phase{}
	// employees := []viewpoint.Employee{}
	//equipment := []viewpoint.

	tx, err := s.epilogue.DB.Beginx()
	s.logger.Info("starting sync...")
	if err != nil {
		s.logger.Error("could not start sync transaction", "err", err.Error())
		return 1
	}

	defer tx.Commit()

	//set sync tables
	file, err := os.ReadFile("internal/db/sql/sync_set_tables.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}
	commands := strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_set_tables.sql")
			tx.Rollback()
			return 1
		}
	}

	//get craftTemplates
	templates := []viewpoint.CraftTemplate{}
	err = s.viewpoint.DB.Select(&templates, syncAllCraftTemplates)
	if err != nil {
		s.logger.Error("failed to get vp craft templates", "err", err.Error())
		tx.Rollback()
		return 1
	}

	_, err = tx.NamedExec(insertTempTableTemplates, templates)
	if err != nil {
		s.logger.Error("failed to insert templates", "err", err.Error())
		tx.Rollback()
		return 1
	}

	file, err = os.ReadFile("internal/db/sql/sync_update_craft_templates.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}

	commands = strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_update_craft_templates.sql")
			tx.Rollback()
			return 1
		}
	}

	//get jobs
	jobs := []viewpoint.Job{}
	err = s.viewpoint.DB.Select(&jobs, syncAllJobs)
	if err != nil {
		s.logger.Error("failed to get vp jobs", "err", err.Error())
		return 1
	}
	_, err = tx.NamedExec(insertTempTableJobs, jobs)
	if err != nil {
		s.logger.Error("failed to insert jobs", "err", err.Error())
		tx.Rollback()
		return 1
	}

	file, err = os.ReadFile("internal/db/sql/sync_update_jobs.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}
	commands = strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_update_jobs.sql")
			err = tx.Rollback()
			if err != nil {
				s.logger.Error("failed to roleback", "err", err.Error())
			}
			return 1
		}
	}

	//get phases
	phases := []viewpoint.Phase{}
	err = s.viewpoint.DB.Select(&phases, syncAllPhases)
	if err != nil {
		s.logger.Error("failed to get vp phases", "err", err.Error())
		return 1
	}
	_, err = tx.NamedExec(insertTempTablePhases, phases)
	if err != nil {
		s.logger.Error("failed to insert phases", "err", err.Error())
		tx.Rollback()
		return 1
	}

	file, err = os.ReadFile("internal/db/sql/sync_update_phases.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}
	commands = strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_update_phases.sql")
			err = tx.Rollback()
			if err != nil {
				s.logger.Error("failed to roleback", "err", err.Error())
			}
			return 1
		}
	}

	//get employees
	employees := []viewpoint.Employee{}
	err = s.viewpoint.DB.Select(&employees, syncAllEmployees)
	if err != nil {
		s.logger.Error("failed to get vp employees", "err", err.Error())
		return 1
	}
	_, err = tx.NamedExec(insertTempTableEmployees, employees)
	if err != nil {
		s.logger.Error("failed to insert employees", "err", err.Error())
		tx.Rollback()
		return 1
	}

	file, err = os.ReadFile("internal/db/sql/sync_update_employees.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}
	commands = strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_update_employees.sql")
			err = tx.Rollback()
			if err != nil {
				s.logger.Error("failed to roleback", "err", err.Error())
			}
			return 1
		}
	}

	//get equipment
	equipment := []viewpoint.Equipment{}
	err = s.viewpoint.DB.Select(&equipment, syncAllEquipment)
	if err != nil {
		s.logger.Error("failed to get vp equipment", "err", err.Error())
		return 1
	}
	_, err = tx.NamedExec(insertTempTableEquipment, equipment)
	if err != nil {
		s.logger.Error("failed to insert equipment", "err", err.Error())
		tx.Rollback()
		return 1
	}

	file, err = os.ReadFile("internal/db/sql/sync_update_equipment.sql")
	if err != nil {
		s.logger.Error("failed to read file", "err", err.Error())
		tx.Rollback()
		return 1
	}
	commands = strings.Split(string(file), ";\n")
	for _, command := range commands {
		_, err = tx.Exec(command)
		if err != nil {
			s.logger.Error("failed to execute sql statement", "err", err.Error(), "file", "sync_update_phases.sql")
			err = tx.Rollback()
			if err != nil {
				s.logger.Error("failed to roleback", "err", err.Error())
			}
			return 1
		}
	}

	tx.Commit()
	s.logger.Info("sync successful... i think")
	return 0
}
