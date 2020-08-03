package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

type RefreshModel struct {
	ID            int64     `db:"id"`
	CreatedAt     time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	TokenRefresh  string    `db:"token_refresh"`
	TokenAccess   string    `db:"token_access"`
	ExpireRefresh time.Time `db:"refresh_expire" time_format:"sql_datetime" time_location:"Local"`
	ExpireAccess  time.Time `db:"access_expire" time_format:"sql_datetime" time_location:"Local"`
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

func GetRefreshFromDb(refreshToken string) RefreshModel {
	a := RefreshModel{}


	return a
}

func (p *RefreshModel) GenerateAccess() AccessModel {
	acc := AccessModel{}
	acc.CreatedAt = time.Now()
	acc.ExpireAt = time.Now().Add(time.Duration(2) * time.Hour)

	acc.UserId = p.UserId
	acc.Scope = p.Scope
	acc.Token = generateAccessToken()

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

	acc := p.GenerateAccess()

	p.ExpireAccess = acc.ExpireAt
	p.TokenAccess = acc.Token

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
	refreshRes, err := db.Exec(
		"insert into `token_refresh`(created_at,token_refresh,token_access,expire_refresh,expire_access,user_id,scope) values(?,?,?,?,?,?,?)",
		refreshModel.CreatedAt,
		refreshModel.TokenRefresh,
		refreshModel.TokenAccess,
		refreshModel.ExpireRefresh,
		refreshModel.ExpireAccess,
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
		"insert into `token_access`(created_at,token,user_id,scope,expire_at) values(?,?,?,?,?)",
		accessModel.CreatedAt,
		accessModel.Token,
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

func (p *RefreshModel) RefreshAccessToken() AccessModel {
	m := NewAccessTokenModel()
	m.Scope = p.Scope
	m.UserId = p.UserId

	//

	oldToken := p.TokenAccess
	oldExpire := p.ExpireAccess
	if oldExpire.After(time.Now()) {
		// expire old token to oldExpire.Add(1 * minute)
		fmt.Println(oldToken)
	}

	return m
}

// RefreshAccessToken 刷新时，如果旧的token还没到期，则到期时间设为 time.Now().Add(1 minute)，即1分钟内，两个token都可以用
func RefreshAccessToken(refreshToken string) {

}
