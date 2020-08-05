package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/oauth/token"
)

func Middleware(c *gin.Context) {
	// token
	query := queryParams{}
	jsonRes := controller.NewJsonRes()

	err := c.ShouldBindQuery(&query)

	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(controller.GetValidationError(err)).Generate())
		c.Abort()
		return
	}

	// get from db
	acc, err := token.GetAccessAvailableWithAppidFromDb(query.AccessToken, query.Appid)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		c.Abort()
		return
	}

	c.Set("accessToken", acc)
	c.Next()
}

type queryParams struct {
	Appid       string `form:"appid" binding:"required,SnowflakeInt64"`
	AccessToken string `form:"access_token" binding:"required,alphanum,len=40"`
}
