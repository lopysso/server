package token

import "github.com/lopysso/server/dependency_injection"

type Generator struct {

	accessToken string

	refreshToken string
}

func (p *Generator) GetAccessToken() string  {
	return p.accessToken
}

func (p *Generator) GetRefreshToken() string {
	return p.refreshToken
}

func (p *Generator) GenerateAccessToken()  {
	
}

func New() *Generator {
	a := Generator{}

	c:= dependency_injection.InjectSnowflakeNode()

	c.Generate().Base32()


	return &a
}


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