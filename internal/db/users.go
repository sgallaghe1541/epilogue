package db

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID              int    `db:"userid"`
	Name            string `db:"name"`
	Email           string `db:"email"`
	Division        string `db:"division"`
	PermissionLevel int    `db:"permissionlevel"`
}

type UserModel struct {
	DB *sqlx.DB
}

func (u *UserModel) ValidEmail(email string) bool {
	var emailresult string
	err := u.DB.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&emailresult)
	return err == nil
}

func (u *UserModel) GetUser(email string) (*User, error) {
	var user = User{}
	err := u.DB.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
