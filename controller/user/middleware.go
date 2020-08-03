package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/libs/account"
)

func Middleware(c *gin.Context) {
	log.Println("hello oauth start")
	token, err := c.Cookie("authorize_token")
	if err != nil {
		c.Redirect(302, "/")
		c.Abort()
		return
	}

	// token
	mo, err := account.GetFromSessionToken(token)
	if err != nil {
		c.Redirect(302, "/")
		c.Abort()
		return
	}

	c.Set("userInfo", mo)

	// find user info

	c.Next()
	log.Println("hello oauth end")
}
