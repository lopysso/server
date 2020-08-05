package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeAction(c *gin.Context) {
	u, err := c.Get("userInfo")
	if !err {
		c.String(http.StatusOK, "user page but error")
		return
	}
	c.String(http.StatusOK, "this is user home \r\n %+v", u)

}
