package code

import (
	"errors"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

const tableName = "oauth_code"

type Model struct {
	Code        string    `db:"code"`
	CreatedAt   time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	Appid       string    `db:"appid"`
	RedirectUri string    `db:"redirect_uri"`
	ExpireAt    time.Time `db:"expire_at" time_format:"sql_datetime" time_location:"LOCAL"`
	Scope       string    `db:"scope"`
	UserId      int64     `db:"user_id"`
}

func NewModelDefault() *Model {
	model := new(Model)

	model.CreatedAt = time.Now()
	// 5分钟
	model.ExpireAt = model.CreatedAt.Add(time.Duration(5) * time.Minute)

	// 先用snowflake 生成code ，但生产环境不能用这个，容易被猜出来
	model.Code = dependency_injection.InjectSnowflakeNode().Generate().Base58()

	return model
}

func (p *Model) Insert() error {

	db := dependency_injection.InjectMysql()

	sqlString := "insert into `" + tableName + "`(`code`,`created_at`,`appid`,`redirect_uri`,`expire_at`,`scope`,`user_id`)"
	sqlString += " values(?,?,?,?,?,?,?)"
	_, err := db.Exec(sqlString, p.Code, p.CreatedAt, p.Appid, p.RedirectUri, p.ExpireAt, p.Scope, p.UserId)

	return err
}

func GetAndDelete(code string) (*Model, error) {
	model := Model{}
	db := dependency_injection.InjectMysql()
	codeRow := db.QueryRow("select code,created_at,appid,redirect_uri,expire_at,scope,user_id from `"+tableName+"` where code=?", code)
	err := codeRow.Scan(
		&model.Code,
		&model.CreatedAt,
		&model.Appid,
		&model.RedirectUri,
		&model.ExpireAt,
		&model.Scope,
		&model.UserId,
	)
	// delete
	defer (func() {
		// 只要有code ，则删除 code
		if len(model.Code) == 0 {
			return
		}
		db := dependency_injection.InjectMysql()
		db.Exec("delete from `"+tableName+"` where code=?", model.Code)
	})()

	if err != nil {
		return nil, errors.New("no ths code")
	}

	if model.ExpireAt.Before(time.Now()) {
		return nil, errors.New("auth code expire")
	}

	return &model, nil
}
