package account

import (
	"errors"
	"fmt"
	"strings"

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
	ID       int64
	Username string
	Password string
	Salt     string
	Status   string
}

// CreateToken 生成token
//
// 暂时这么搞
func (p *Model) CreateToken() string {

	return fmt.Sprintf("%s::%s", p.Username, p.Password)
}

// GetFromDb 从数据库取model
func GetFromDb(username string) (*Model, error) {
	mo := Model{}
	db := dependency_injection.InjectMysql()

	// var rowUsername,rowPassword,rowID,rowSalt string
	userRow := db.QueryRow("select id,username,password,salt,status from user where username=?", username)

	err := userRow.Scan(&mo.ID, &mo.Username, &mo.Password, &mo.Salt, &mo.Status)
	// err := userRow.Scan(mo.ID,mo.Username,mo.Password,mo.Salt,mo.Status)
	// err := userRow.Scan(mo)
	if err != nil {
		return nil, err
	}

	return &mo, err
}

// GetWithPassword 从数据库中获取model
func GetWithPassword(username string, password string) (*Model, error) {

	mo, err := GetFromDb(username)
	if err != nil {
		return nil, err
	}

	pwdHash := HashPwd(password, mo.Salt)
	if pwdHash != mo.Password {
		return nil, errors.New("password error")
	}

	return mo, nil
}

// GetFromToken 从token 中取model
func GetFromToken(token string) (*Model, error) {

	userArr := strings.Split(token, "::")
	if len(userArr) != 2 {
		return nil, errors.New("token error")
	}

	mo, err := GetFromDb(userArr[0])
	if err != nil {
		return nil, err
	}

	// check
	if mo.Password != userArr[1] {
		return nil, errors.New("token invalid")
	}

	return mo, nil
}
