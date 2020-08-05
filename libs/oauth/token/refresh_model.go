package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

// 宽限时间
// 即即使超时了，还是可以用一下， 这里看看怎么设计
const expireGrace = 60

type RefreshModel struct {
	ID            int64     `db:"id"`
	CreatedAt     time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	TokenRefresh  string    `db:"token_refresh"`
	TokenAccess   string    `db:"token_access"`
	ExpireRefresh time.Time `db:"refresh_expire" time_format:"sql_datetime" time_location:"Local"`
	ExpireAccess  time.Time `db:"access_expire" time_format:"sql_datetime" time_location:"Local"`
	Appid         string    `db:"appid"`
	UserId        int64     `db:"user_id"`
	Scope         string    `db:"scope"`
}

func NewRefresh() RefreshModel {
	a := RefreshModel{}
	a.CreatedAt = time.Now()
	a.ExpireRefresh = time.Now().AddDate(0, 0, 30)
	a.TokenRefresh = generateRefreshToken()

	return a
}

func GetRefreshFromDb(refreshToken string) (*RefreshModel, error) {
	a := RefreshModel{}
	db := dependency_injection.InjectMysql()

	cols := "id,created_at,token_refresh,token_access,expire_refresh,expire_access,appid,user_id,scope"
	sqlString := fmt.Sprintf("select %s from token_refresh where token_refresh=?", cols)
	row := db.QueryRow(sqlString, refreshToken)
	err := row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.TokenRefresh,
		&a.TokenAccess,
		&a.ExpireRefresh,
		&a.ExpireAccess,
		&a.Appid,
		&a.UserId,
		&a.Scope,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func GetRefreshWithAppidFromDb(refreshToken string, appid string) (*RefreshModel, error) {
	a, err := GetRefreshFromDb(refreshToken)
	if err != nil {
		return nil, err
	}

	if a.Appid != appid {
		return nil, errors.New("refreshToken can not matched")
	}

	return a, nil
}

func (p *RefreshModel) GenerateAccess() AccessModel {
	acc := AccessModel{}
	acc.CreatedAt = time.Now()
	acc.ExpireAt = time.Now().Add(time.Duration(2) * time.Hour)

	acc.Appid = p.Appid
	acc.UserId = p.UserId
	acc.Scope = p.Scope
	acc.Token = generateAccessToken()

	return acc
}

func (p *RefreshModel) generateAccessAndUpdateSelf() AccessModel {
	acc := p.GenerateAccess()

	p.ExpireAccess = acc.ExpireAt
	p.TokenAccess = acc.Token

	return acc
}

func (p *RefreshModel) InsertToDb() (*AccessModel, error) {
	if p.ID > 0 {
		return nil, errors.New("can not insert to table")
	}

	if p.UserId <= 0 {
		return nil, errors.New("has not set user")
	}

	// new access token

	// acc := p.GenerateAccess()

	// p.ExpireAccess = acc.ExpireAt
	// p.TokenAccess = acc.Token
	acc := p.generateAccessAndUpdateSelf()

	err := insertTokens(p, &acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func insertTokens(refreshModel *RefreshModel, accessModel *AccessModel) error {
	db := dependency_injection.InjectMysql()
	trans, err := db.Begin()
	if err != nil {
		return errors.New("db error")
	}

	// refresh token
	sqlString := "insert into `token_refresh`(created_at,token_refresh,token_access,expire_refresh,expire_access,appid,user_id,scope) "
	sqlString += " values(?,?,?,?,?,?,?,?)"
	refreshRes, err := db.Exec(
		sqlString,
		refreshModel.CreatedAt,
		refreshModel.TokenRefresh,
		refreshModel.TokenAccess,
		refreshModel.ExpireRefresh,
		refreshModel.ExpireAccess,
		refreshModel.Appid,
		refreshModel.UserId,
		refreshModel.Scope,
	)
	if err != nil {
		trans.Rollback()
		return err
	}

	refreshModel.ID, _ = refreshRes.LastInsertId()

	// access token
	accessRes, err := db.Exec(
		"insert into `token_access`(created_at,token,appid,user_id,scope,expire_at) values(?,?,?,?,?,?)",
		accessModel.CreatedAt,
		accessModel.Token,
		accessModel.Appid,
		accessModel.UserId,
		accessModel.Scope,
		accessModel.ExpireAt,
	)

	if err != nil {
		trans.Rollback()
		return err
	}
	accessModel.ID, _ = accessRes.LastInsertId()

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return err
	}

	return nil
}

// 刷新token
// 如果旧还没有超时，则将旧的改超时
func (p *RefreshModel) RefreshAccessToken() (*AccessModel, error) {
	// m := NewAccessTokenModel()
	// m.Scope = p.Scope
	// m.UserId = p.UserId
	// m.Appid = p.Appid

	//

	db := dependency_injection.InjectMysql()
	trans, err := db.Begin()
	if err != nil {

		return nil, err
	}

	if p.ExpireAccess.After(time.Now()) {
		// 强制超时
		newExpire := time.Now().Add(time.Duration(expireGrace) * time.Second)
		_, err := db.Exec("update token_access set expire_at=? where token=?", newExpire, p.TokenAccess)
		if err != nil {
			log.Println(err)
		}
	}

	// new token
	accessModel := p.generateAccessAndUpdateSelf()
	_, err = db.Exec("update token_refresh set token_access=?,expire_access=?", p.TokenAccess, p.ExpireAccess)
	if err != nil {
		trans.Rollback()
		return nil, err
	}

	accessRes, err := db.Exec(
		"insert into `token_access`(created_at,token,appid,user_id,scope,expire_at) values(?,?,?,?,?,?)",
		accessModel.CreatedAt,
		accessModel.Token,
		accessModel.Appid,
		accessModel.UserId,
		accessModel.Scope,
		accessModel.ExpireAt,
	)

	accessModel.ID, _ = accessRes.LastInsertId()

	if err != nil {
		trans.Rollback()
		return nil, err
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return nil, err
	}
	// refresh to db

	return &accessModel, nil
}
