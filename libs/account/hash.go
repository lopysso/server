package account

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"

	"github.com/lopysso/server/libs/text/random"
)

// var base58Alphabets = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func hashSha1(str string) string {

	data := []byte(str)
	has := sha1.Sum(data)
	l := fmt.Sprintf("%x", has)
	return l
}

func hashMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func HashPwd(password string, salt string) string {
	return HashPwdFromMd5(hashMd5(password), salt)
}

func HashPwdFromMd5(passwordMd5 string, salt string) string {

	s := hashMd5(passwordMd5+salt) + hashMd5(salt)

	return hashSha1(s)
}

// CreateSalt 找一个随机算法
func CreateSalt(length int) string {

	return random.Base58(length)
}
