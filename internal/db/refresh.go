package db

import (
	"errors"
	"time"
)

var ErrRefreshExpired = errors.New("refresh token expired")

type RefreshToken struct {
	UserID  int       `db:"userid"`
	Token   string    `db:"token"`
	Expires time.Time `db:"expires"`
}

func (e *EpilogueConnection) SaveRefreshToken(userID int, token string) error {
	refreshToken := RefreshToken{
		UserID:  userID,
		Token:   token,
		Expires: time.Now().Add(time.Hour),
	}

	_, err := e.DB.NamedExec("INSERT INTO refreshtokens (userid, token, expires) VALUES (:userid, :token, :expires)", refreshToken)

	return err
}

func (e *EpilogueConnection) UserForRefreshToken(token string) (int, error) {
	refreshToken := RefreshToken{}

	err := e.DB.Select(&refreshToken, "SELECT * FROM refreshtokens WHERE token = ?", token)
	if err != nil {
		return 0, err
	}

	if refreshToken.Expires.Before(time.Now()) {
		return 0, ErrRefreshExpired
	}

	return refreshToken.UserID, nil
}

func (e *EpilogueConnection) RevokeRefreshToken(token string) error {
	_, err := e.DB.Exec("DELETE FROM refreshtokens WHERE token = ?", token)
	return err
}
