package user

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/account"
	"github.com/lopysso/server/libs/oauth/app"
	"github.com/lopysso/server/libs/oauth/code"
)

// "fmt"
// "log"
// "net/http"
// "net/url"

// "github.com/gin-gonic/gin"
// "github.com/go-playground/validator/v10"
// "github.com/lopysso/server/libs/account"
// "github.com/lopysso/server/libs/oauth/app"
// "github.com/lopysso/server/libs/oauth/code"

// Authorize 用创建授权码code
//
// 这个是不是得改
// 正常情况下，可以直接302，也可以让用户点是否授权，就像github那样
// 没登录则跳到登录页面
func Authorize(c *gin.Context) {

	query := oauthQueryParams{}
	err := c.ShouldBindQuery(&query)
	if err != nil {

		c.String(http.StatusOK, "authorize error, %+v", controller.GetValidationError(err))
		return
	}

	log.Printf("query is : %+v", query)

	u, err := url.Parse(query.RedirectUri)
	if err != nil {
		c.String(http.StatusOK, "authorize error, %+v", err)
		return
	}

	// find appid
	appRow, err := app.GetFromAppid(query.Appid)
	if err != nil {
		c.String(http.StatusOK, "can not find this app")
		return
	}

	log.Printf("%+v", appRow)

	// get userinfo

	userInfo, exist := c.Get("userInfo")
	if !exist {
		c.String(http.StatusOK, "user not login")
		return
	}
	log.Printf("%+v", userInfo)
	userInfoNew := userInfo.(*account.Model)
	log.Printf("%+v", userInfoNew)

	// insert code
	codeModel := code.NewModelDefault()
	codeModel.Appid = appRow.Appid
	codeModel.RedirectUri = query.RedirectUri
	codeModel.Scope = query.Scope
	codeModel.UserId = userInfoNew.ID
	err = codeModel.Insert()
	if err != nil {
		c.String(http.StatusOK, "insert code error %+v", err)
		return
	}

	//
	uQuery := u.Query()

	uQuery.Set("code", codeModel.Code)
	uQuery.Set("state", query.State)
	u.RawQuery = uQuery.Encode()

	// 哪种情况下，直接redirect，哪种情况下，让用户点击是否授权
	// 现在，直接redirect
	log.Printf("redirect to : %s", u.String())
	// c.String(http.StatusOK, "to do redirect: %s", u.String())
	c.Redirect(http.StatusMovedPermanently, u.String())
}

// oauthQueryParams
//
// 后续自己写一些验证器
type oauthQueryParams struct {
	Appid        string `form:"appid" binding:"required,SnowflakeInt64"`
	ResponseType string `form:"response_type" binding:"required"`
	RedirectUri  string `form:"redirect_uri" binding:"required"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
}
