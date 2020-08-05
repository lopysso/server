package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/oauth/token"
)

// 用于刷新access_token
func RefreshTokenAction(c *gin.Context) {
	query := refreshQueryParams{}
	jsonRes := controller.NewJsonRes()
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(controller.GetValidationError(err)).Generate())
		return
	}

	refreshModel, err := token.GetRefreshWithAppidFromDb(query.RefreshToken, query.Appid)
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	accessTokenModle, err := refreshModel.RefreshAccessToken()
	if err != nil {
		c.JSON(http.StatusOK, jsonRes.Error(err.Error()).Generate())
		return
	}

	jsonRes.Data(gin.H{
		"access_token":  refreshModel.TokenAccess,
		"token_type":    "",
		"expires_in":    accessTokenModle.GetExpireIn(),
		"refresh_token": refreshModel.TokenRefresh,
	})

	c.JSON(http.StatusOK, jsonRes.Success().Generate())
}

type refreshQueryParams struct {
	Appid        string `form:"appid" binding:"required,SnowflakeInt64"`
	RefreshToken string `form:"refresh_token" binding:"required,alphanum,min=48,max=48"`
}
