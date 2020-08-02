package token

import (
	"errors"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

type AccessModel struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	Token     string    `db:"token"`
	UserId    int64     `db:"user_id"`
	Scope     string    `db:"scope"`
	ExpireAt  time.Time `db:"expire_at" time_format:"sql_datetime" time_location:"Local"`
}

func NewAccessTokenModel() AccessModel {
	a := AccessModel{}
	a.CreatedAt = time.Now()
	a.ExpireAt = time.Now().Add(time.Duration(2) * time.Hour)
	return a
}

type RefreshModel struct {
	ID           int64     `db:"id"`
	CreatedAt    time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	TokenRefresh string    `db:"token_refresh"`
	TokenAccess  string    `db:"token_access"`
	ExpireAt     time.Time `db:"expire_at" time_format:"sql_datetime" time_location:"Local"`
}

func NewRefreshTokenModel() RefreshModel {
	a := RefreshModel{}
	a.CreatedAt = time.Now()
	a.ExpireAt = time.Now().AddDate(0, 0, 30)
	return a
}

func InsertTokens(accessModel *AccessModel, refreshModel *RefreshModel) error {
	db := dependency_injection.InjectMysql()
	trans, err := db.Begin()
	if err != nil {
		return errors.New("db error")
	}

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

	// refresh token
	refreshRes, err := db.Exec(
		"insert into `token_refresh`(created_at,token_refresh,token_access,expire_at) values(?,?,?,?)",
		refreshModel.CreatedAt,
		refreshModel.TokenRefresh,
		refreshModel.TokenAccess,
		refreshModel.ExpireAt,
	)
	if err != nil {
		trans.Rollback()
		return err
	}

	refreshModel.ID, _ = refreshRes.LastInsertId()

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return err
	}

	return nil
}

// RefreshAccessToken 刷新时，如果旧的token还没到期，则到期时间设为 time.Now().Add(1 minute)，即1分钟内，两个token都可以用
func RefreshAccessToken(refreshToken string) {

}
