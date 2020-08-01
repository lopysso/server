package index

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	// "github.com/lopysso/server/dependency_injection"
	"github.com/lopysso/server/libs/account"
)

// SigninAction 接收post 登录
func SigninAction(c *gin.Context) {
	jsonRes := gin.H{
		"http": "true",
		"code": "1",
		"msg":  "",
		"data": nil,
	}

	user := userForm{}

	err := c.ShouldBind(&user)

	log.Printf("%+v\n", user)

	if err != nil {
		jsonRes["msg"] = user.getError(err.(validator.ValidationErrors))
		c.JSON(http.StatusOK, jsonRes)
		return
	}

	log.Printf("form data success")

	jsonRes["data"] = gin.H{
		"username": user.Username,
	}

	userModel, err := account.GetWithPassword(user.Username, user.Password)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, jsonRes)
		return
	}

	log.Println("login ok")
	log.Println(userModel)

	// 暂时不用token， 先这么搞 username::password
	c.SetCookie("authorize_token", userModel.CreateToken(), 99999999, "/", "", false, true)

	jsonRes["code"] = "0"
	jsonRes["msg"] = "login ok"
	c.JSON(http.StatusOK, jsonRes)
}

type userForm struct {
	Username string `form:"username" binding:"required,alphanum,min=4,max=16"`
	Password string `form:"password" binding:"required,min=6,max=16"`
}

func (r *userForm) getError(err validator.ValidationErrors) string {

	for _, v := range err {
		log.Println(v.Field(), v.Tag())
		//return "error test"
		errMsg := "格式错误"
		if v.Tag() == "required" {
			errMsg = "不能为空"
		}
		return fmt.Sprintf("%s %s", v.Field(), errMsg)
	}

	return "unknown error"
}
