package token

import (
	"errors"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

const accessExpireInSecond = 7200

type AccessModel struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	Token     string    `db:"token"`
	Appid     string    `db:"appid"`
	UserId    int64     `db:"user_id"`
	Scope     string    `db:"scope"`
	ExpireAt  time.Time `db:"expire_at" time_format:"sql_datetime" time_location:"Local"`
}

func NewAccessTokenModel() AccessModel {
	a := AccessModel{}
	a.CreatedAt = time.Now()
	a.ExpireAt = time.Now().Add(time.Duration(accessExpireInSecond) * time.Second)
	return a
}

func (p *AccessModel) GetExpireIn() int {
	return accessExpireInSecond
}

func GetAccessAvailableFromDb(token string) (*AccessModel, error) {
	m := AccessModel{}
	db := dependency_injection.InjectMysql()
	row := db.QueryRow("select id,created_at,token,appid,user_id,scope,expire_at from token_access where token=?", token)
	err := row.Scan(
		&m.ID,
		&m.CreatedAt,
		&m.Token,
		&m.Appid,
		&m.UserId,
		&m.Scope,
		&m.ExpireAt,
	)
	if err != nil {
		return nil, err
	}

	// expire
	if m.ExpireAt.Before(time.Now()) {
		return nil, errors.New("token expire")
	}

	return &m, nil
}

func GetAccessAvailableWithAppidFromDb(token string, appid string) (*AccessModel, error) {
	acc, err := GetAccessAvailableFromDb(token)
	if err != nil {
		return nil, err
	}

	if acc.Appid != appid {
		return nil, errors.New("appid has not matched this token")
	}

	return acc, nil
}
