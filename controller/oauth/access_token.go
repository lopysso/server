package oauth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lopysso/server/libs/oauth/code"
	"github.com/lopysso/server/libs/oauth/token"
)

func AccessToken(c *gin.Context) {

	jsonRes := gin.H{
		"http": "true",
		"code": "1",
		"msg":  "",
		"data": nil,
	}

	query := queryParams{}
	err := c.ShouldBindQuery(&query)

	if err != nil {
		jsonRes["msg"] = query.getError(err.(validator.ValidationErrors))
		c.JSON(http.StatusOK, jsonRes)
		return
	}

	// get and delete code, trans
	codeModel, err := code.GetAndDelete(query.Code)
	if err != nil {
		jsonRes["msg"] = err.Error()
		c.JSON(http.StatusOK, jsonRes)
		return
	}

	log.Printf("code model : %+v", codeModel)

	// insert access_token and refresh_token
	tokenObj := token.New()

	log.Println(tokenObj.GetRefreshToken())
	log.Println(tokenObj.GetRefreshToken())
	log.Println(tokenObj.GetAccessToken())
	log.Println(tokenObj.GetAccessToken())
	log.Println(tokenObj.GenerateAccessToken())
	log.Println(tokenObj.GetAccessToken())
	log.Println(tokenObj.GetAccessToken())

	accessModel := token.NewAccessTokenModel()
	accessModel.Token = tokenObj.GetAccessToken()
	accessModel.UserId = codeModel.UserId
	accessModel.Scope = codeModel.Scope

	refreshModel := token.NewRefreshTokenModel()
	refreshModel.TokenAccess = tokenObj.GetAccessToken()
	refreshModel.TokenRefresh = tokenObj.GetRefreshToken()

	err = token.InsertTokens(&accessModel, &refreshModel)
	if err != nil {
		jsonRes["msg"] = err.Error()
		c.JSON(http.StatusOK, jsonRes)
		return
	}

	//

	jsonRes["code"] = "0"
	jsonRes["msg"] = "ok"
	jsonRes["data"] = gin.H{
		"access_token":  refreshModel.TokenAccess,
		"token_type":    "",
		"expires_in":    7200,
		"refresh_token": refreshModel.TokenRefresh,
	}

	c.JSON(200, jsonRes)

	// c.String(http.StatusOK, "hehe %+v",query)
}

const AuthorizationCode = "authorization_code"

type queryParams struct {
	Appid     string `form:"appid" binding:"required,alphanum"`
	Secret    string `form:"secret" binding:"required,alphanum"`
	Code      string `form:"code" binding:"required,alphanum"`
	GrantType string `form:"grant_type" binding:"required"`
}

func (r *queryParams) getError(err validator.ValidationErrors) string {

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
