package oauth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/oauth/app"
	"github.com/lopysso/server/libs/oauth/code"
	"github.com/lopysso/server/libs/oauth/token"
)

func AccessToken(c *gin.Context) {

	jsonRes := controller.NewJsonRes()

	query := accessQueryParams{}
	err := c.ShouldBindQuery(&query)

	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(controller.GetValidationError(err)).Generate())
		return
	}

	// check appid and secret
	_, err = app.GetFromAppidWithSecret(query.Appid, query.Secret)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error("can not use this code").Generate())
		return
	}

	// get and delete code, trans
	codeModel, err := code.GetAndDelete(query.Code)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	log.Printf("code model : %+v", codeModel)

	refreshModel := token.NewRefresh()
	refreshModel.Scope = codeModel.Scope
	refreshModel.UserId = codeModel.UserId
	refreshModel.Appid = codeModel.Appid

	accessTokenModle, err := refreshModel.InsertToDb()
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	//

	jsonRes.Success().Data(gin.H{
		"access_token":  refreshModel.TokenAccess,
		"token_type":    "",
		"expires_in":    accessTokenModle.GetExpireIn(),
		"refresh_token": refreshModel.TokenRefresh,
	})
	c.JSON(200, jsonRes.Generate())

	// c.String(http.StatusOK, "hehe %+v",query)
}

const AuthorizationCode = "authorization_code"

type accessQueryParams struct {
	Appid     string `form:"appid" binding:"required,SnowflakeInt64"`
	Secret    string `form:"secret" binding:"required,alphanum"`
	Code      string `form:"code" binding:"required,alphanum"`
	GrantType string `form:"grant_type" binding:"required"`
}

func (p *accessQueryParams) AppidInt() int64 {
	a, _ := strconv.ParseInt(p.Appid, 10, 64)

	return a
}
