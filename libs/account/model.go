package account

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lopysso/server/dependency_injection"
)

// StatusNormal 正常状态
const StatusNormal = "NORMAL"

// StatusClose 关闭状态
const StatusClose = "CLOSE"

// StatusDelete 删除状态
const StatusDelete = "DELETE"

// Model 用户信息
type Model struct {
	ID        int64
	CreatedAt time.Time `db:"created_at" time_format:"sql_datetime" time_location:"Local"`
	Nickname  string
	Username  string
	Password  string
	Salt      string
	Status    string
}

// CreateSessionToken 生成token
//
// 暂时这么搞
func (p *Model) CreateSessionToken() string {

	return fmt.Sprintf("%s::%s", p.Username, p.Password)
}

// GetFromUsername 从数据库取model
func GetFromID(id int64) (*Model, error) {
	mo := Model{}
	db := dependency_injection.InjectMysql()

	// var rowUsername,rowPassword,rowID,rowSalt string
	userRow := db.QueryRow("select id,created_at,nickname,username,password,salt,status from user where id=?", id)

	err := userRow.Scan(&mo.ID, &mo.CreatedAt, &mo.Nickname, &mo.Username, &mo.Password, &mo.Salt, &mo.Status)
	if err != nil {
		return nil, err
	}

	return &mo, err
}

// GetFromUsername 从数据库取model
func GetFromUsername(username string) (*Model, error) {
	mo := Model{}
	db := dependency_injection.InjectMysql()

	// var rowUsername,rowPassword,rowID,rowSalt string
	userRow := db.QueryRow("select id,created_at,nickname,username,password,salt,status from user where username=?", username)

	err := userRow.Scan(&mo.ID, &mo.CreatedAt, &mo.Nickname, &mo.Username, &mo.Password, &mo.Salt, &mo.Status)
	// err := userRow.Scan(mo.ID,mo.Username,mo.Password,mo.Salt,mo.Status)
	// err := userRow.Scan(mo)
	if err != nil {
		return nil, err
	}

	return &mo, err
}

// GetWithPassword 从数据库中获取model
func GetWithPassword(username string, password string) (*Model, error) {

	mo, err := GetFromUsername(username)
	if err != nil {
		return nil, err
	}

	pwdHash := HashPwd(password, mo.Salt)
	if pwdHash != mo.Password {
		return nil, errors.New("password error")
	}

	return mo, nil
}

// GetFromSessionToken 从session 中取model
func GetFromSessionToken(token string) (*Model, error) {

	userArr := strings.Split(token, "::")
	if len(userArr) != 2 {
		return nil, errors.New("token error")
	}

	mo, err := GetFromUsername(userArr[0])
	if err != nil {
		return nil, err
	}

	// check
	if mo.Password != userArr[1] {
		return nil, errors.New("token invalid")
	}

	return mo, nil
}

func (p *Model) ChangePassword(oldPassword string, newPassword string) error {
	if p.ID <= 0 {
		return errors.New("server error: no this user")
	}
	if HashPwd(oldPassword, p.Salt) != p.Password {
		return errors.New("old password error")
	}

	//
	newSalt := CreateSalt(8)
	newHash := HashPwd(newPassword, newSalt)

	db := dependency_injection.InjectMysql()
	_, err := db.Exec("update user set password=?,salt=? where id=?", newHash, newSalt, p.ID)
	if err != nil {
		return errors.New("password changed error")
	}

	p.Password = newHash
	p.Salt = newSalt

	return nil

}
