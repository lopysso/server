package token

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lopysso/server/dependency_injection"
)

func generateRefreshToken() string {
	flakeNode := dependency_injection.InjectSnowflakeNode()

	num := flakeNode.Generate()

	str := fmt.Sprintf("%s::%d", num.Base32(), num.Int64())

	return hashRefreshToken(str)
}

func generateAccessToken() string {
	flakeNode := dependency_injection.InjectSnowflakeNode()

	num := flakeNode.Generate()

	str := fmt.Sprintf("%s:::%d", num.Base32(), num.Int64())

	return hashAccessToken(str)
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
