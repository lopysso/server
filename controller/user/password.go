package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/account"
)

func PasswordAction(c *gin.Context) {
	userInfo, _ := c.Get("userInfo")
	userInfo1, ok := userInfo.(*account.Model)

	log.Println(userInfo1, ok)
	c.HTML(http.StatusOK, "user/index.html", gin.H{
		"title":    "hello",
		"userInfo": userInfo1,
	})
}

func PasswordUpdateAction(c *gin.Context) {
	userInfo, _ := c.Get("userInfo")
	userInfo1, ok := userInfo.(*account.Model)
	jsonRes := controller.NewJsonRes()
	if !ok {
		c.JSON(http.StatusOK, jsonRes.Error("server error: userinfo error").Generate())
		return
	}

	postData := passwordPost{}
	err := c.ShouldBind(&postData)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(controller.GetValidationError(err)).Generate())
		return
	}

	err = userInfo1.ChangePassword(postData.Password, postData.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	// re login or write new session

	c.SetCookie("authorize_token", userInfo1.CreateSessionToken(), 99999999, "/", "", false, true)
	c.JSON(http.StatusOK, jsonRes.Success().Generate())
}

type passwordPost struct {
	Password    string `form:"password" binding:"required,min=6,max=20"`
	NewPassword string `form:"new_password" binding:"required,min=6,max=20"`
	RePassword  string `form:"re_password" binding:"required,min=6,max=20,eqfield=NewPassword"`
	Capcha      string
}
