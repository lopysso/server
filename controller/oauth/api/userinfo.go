package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/account"
	"github.com/lopysso/server/libs/oauth/token"
)

func Userinfo(c *gin.Context)  {
	acc ,exists := c.Get("accessToken")
	jsonRes := controller.NewJsonRes()
	if !exists{
		c.JSON(http.StatusOK,jsonRes.Error("server error").Generate())
		return
	}
	accessTokenModel := acc.(*token.AccessModel)

	userInfo ,err := account.GetFromID(accessTokenModel.UserId)
	if err != nil{
		c.JSON(http.StatusOK,jsonRes.Error(err.Error()).Generate())
		return
	}

	data := gin.H{
		"id":userInfo.ID,
		"nickname":userInfo.Nickname,
		"username":userInfo.Username,
		"status":userInfo.Status,
	}

	c.JSON(http.StatusOK,jsonRes.Success().Data(gin.H{
		"user":data,
	}).Generate())
}