package db

import (
	"fmt"
)

type Permission struct {
	Department string `db:"department" json:"department"`
	Role       string `db:"roleid" json:"roleid"`
}

func (p Permission) Stringify() string {
	return fmt.Sprintf("%s|%s,", p.Department, p.Role)
}

type User struct {
	ID          int          `db:"userid"`
	Name        string       `db:"name"`
	Email       string       `db:"email"`
	Active      int          `db:"active"`
	Permissions []Permission `json:"permissions"`
}

func (e *EpilogueConnection) ValidEmail(email string) bool {
	var emailresult string
	err := e.DB.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&emailresult)
	return err == nil
}

func (e *EpilogueConnection) ActiveUser(id int) (bool, error) {
	var active int
	err := e.DB.QueryRow("SELECT active FROM users WHERE userid=?", id).Scan(&active)
	if err != nil {
		return false, ErrUserDoesNotExist
	}
	if active != 1 {
		return false, ErrInactiveUser
	}
	return true, nil
}

func (e *EpilogueConnection) GetUser(email string) (*User, error) {
	var (
		user        = User{}
		permissions = []Permission{}
	)
	err := e.DB.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	if user.Active == 0 {
		return nil, ErrInactiveUser
	}
	err = e.DB.Select(&permissions, "SELECT department, roleid FROM user_permissions WHERE userid=?", user.ID)
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions found for %s", user.Name)
	}

	user.Permissions = permissions

	return &user, nil
}
