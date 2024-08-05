package db

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrRefreshExpired = errors.New("refresh token expired")

type RefreshToken struct {
	UserID  int       `db:"userid"`
	Token   string    `db:"token"`
	Expires time.Time `db:"expires"`
}

type RefreshTokenModel struct {
	DB *sqlx.DB
}

func (r *RefreshTokenModel) SaveRefreshToken(userID int, token string) error {
	refreshToken := RefreshToken{
		UserID:  userID,
		Token:   token,
		Expires: time.Now().Add(time.Hour),
	}

	_, err := r.DB.NamedExec("INSERT INTO refreshtokens (userid, token, expires) VALUES (:userid, :token, :expires)", refreshToken)

	return err
}

func (r *RefreshTokenModel) UserForRefreshToken(token string) (int, error) {
	refreshToken := RefreshToken{}

	err := r.DB.Select(&refreshToken, "SELECT * FROM refreshtokens WHERE token = ?", token)
	if err != nil {
		return 0, err
	}

	if refreshToken.Expires.Before(time.Now()) {
		return 0, ErrRefreshExpired
	}

	return refreshToken.UserID, nil
}

func (r *RefreshTokenModel) RevokeRefreshToken(token string) error {
	_, err := r.DB.Exec("DELETE FROM refreshtokens WHERE token = ?", token)
	return err
}
