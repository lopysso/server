package index

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	// "github.com/lopysso/server/dependency_injection"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/account"
)

// SigninAction 接收post 登录
func SigninAction(c *gin.Context) {

	jsonRes := controller.NewJsonRes()

	user := userForm{}

	err := c.ShouldBind(&user)

	log.Printf("%+v\n", user)

	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(controller.GetValidationError(err)).Generate())
		return
	}

	log.Printf("form data success")

	jsonRes.Data(gin.H{
		"username": user.Username,
	})

	userModel, err := account.GetWithPassword(user.Username, user.Password)

	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	log.Println("login ok")
	log.Println(userModel)

	// 暂时不用token， 先这么搞 username::password
	c.SetCookie("authorize_token", userModel.CreateSessionToken(), 99999999, "/", "", false, true)

	jsonRes.Success("login ok")
	c.JSON(http.StatusOK, jsonRes.Generate())
}

type userForm struct {
	Username string `form:"username" binding:"required,alphanum,min=4,max=20"`
	Password string `form:"password" binding:"required,min=6,max=20"`
}
