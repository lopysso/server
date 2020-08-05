package app

import (
	"errors"
	"log"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

type Model struct {
	Id        int
	CreatedAt time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	Appid     string
	Secret    string
	Title     string
	Desp      string
	Domains   string
}

func GetFromAppid(appid string) (*Model, error) {
	model := Model{}

	a := dependency_injection.InjectMysql().QueryRow("select id,created_at,appid,secret,title,desp,domains from app where appid=?", appid)
	err := a.Scan(
		&model.Id,
		&model.CreatedAt,
		&model.Appid,
		&model.Secret,
		&model.Title,
		&model.Desp,
		&model.Domains,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model, nil
}

func GetFromAppidWithSecret(appid string, secret string) (*Model, error) {
	m, err := GetFromAppid(appid)

	if err != nil {
		return nil, err
	}

	if m.Secret != secret {
		return nil, errors.New("app secret error")
	}

	return m, nil
}
