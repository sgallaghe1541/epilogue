package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Permission struct {
	Division string `db:"division" json:"division"`
	Role     string `db:"role" json:"role"`
}

func (p Permission) Stringify() string {
	return fmt.Sprintf("%s|%s,", p.Division, p.Role)
}

type User struct {
	ID          int          `db:"userid"`
	Name        string       `db:"name"`
	Email       string       `db:"email"`
	Active      int          `db:"active"`
	Permissions []Permission `json:"permissions"`
}

type UserModel struct {
	DB *sqlx.DB
}

func (u *UserModel) ValidEmail(email string) bool {
	var emailresult string
	err := u.DB.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&emailresult)
	return err == nil
}

func (u *UserModel) Active(id int) (bool, error) {
	var active int
	err := u.DB.QueryRow("SELECT active FROM users WHERE userid=?", id).Scan(&active)
	if err != nil {
		return false, ErrUserDoesNotExist
	}
	if active != 1 {
		return false, ErrInactiveUser
	}
	return true, nil
}

func (u *UserModel) GetUser(email string) (*User, error) {
	var (
		user        = User{}
		permissions = []Permission{}
	)
	err := u.DB.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	if user.Active == 0 {
		return nil, ErrInactiveUser
	}
	err = u.DB.Select(&permissions, "SELECT division, role FROM user_permissions WHERE userid=?", user.ID)
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions found for %s", user.Name)
	}

	user.Permissions = permissions

	return &user, nil
}
