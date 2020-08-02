package token

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lopysso/server/dependency_injection"
)

type Generator struct {
	accessToken string

	refreshToken string
}

func (p *Generator) GetAccessToken() string {
	return p.accessToken
}

func (p *Generator) GetRefreshToken() string {
	return p.refreshToken
}

func (p *Generator) GenerateAccessToken() string {
	flakeNode := dependency_injection.InjectSnowflakeNode()

	num := flakeNode.Generate()

	str := fmt.Sprintf("%s:::%d", num.Base32(), num.Int64())

	p.accessToken = hashAccessToken(str)
	return p.accessToken
}

func generateRefreshToken() string {
	flakeNode := dependency_injection.InjectSnowflakeNode()

	num := flakeNode.Generate()

	str := fmt.Sprintf("%s::%d", num.Base32(), num.Int64())

	return hashRefreshToken(str)
}

func New() *Generator {
	a := Generator{}

	a.refreshToken = generateRefreshToken()
	a.GenerateAccessToken()

	return &a
}

func hashRefreshToken(str string) string {
	data := []byte(str)
	// sum := sha256.Sum224(data)
	sum := sha256.Sum256(data)

	res := base58.Encode(sum[:])
	res += randomString(5)

	// 48
	return res[0:48]
}

func hashAccessToken(str string) string {
	data := []byte(str)
	sum := sha256.Sum224(data)
	// sum := sha256.Sum256(data)

	res := base58.Encode(sum[:])
	res += randomString(3)

	// 48
	return res[0:40]
}

func randomString(l int) string {
	str := "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
