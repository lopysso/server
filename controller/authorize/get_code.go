package authorize

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCode 用创建授权码code
//
// 这个是不是得改
// 正常情况下，可以直接302，也可以让用户点是否授权，就像github那样
func GetCode(c *gin.Context) {

	// c.Query("response_type") // 指定为code,目前 没有啥用
	// c.Query("client_id")
	// c.Query("redirect_uri")
	// c.Query("scope") // 权限 ，可选
	// 将client_id redirect_uri code 这几项存在临时数据库里，比如变量？redis?
	// 该三项必须一一对应。然后  code的效期不能太长，一般为几分钟，微信好像是5分钟

	// 至少
	dataMap := make(map[string]interface{})
	dataMap["code"] = "create a random"
	// 原路返回
	dataMap["state"] = c.DefaultQuery("state", "")

	log.Println(dataMap)
	c.String(200, "do redirect to request_uri with code and state or show a page to confirm")
	// c.Redirect(302, "request_uri")
}

// GetToken 从app 的服务端来的，请求token and refresh_token
// 这个是
func GetToken(c *gin.Context) {
	// c.Query("client_id")
	// c.Query("redirect_uri")
	// c.Query("code")
	// client_id redirect_uri code 这几项必须一一对应，

	dataMap := make(map[string]interface{})

	// 表示访问令牌，必有
	dataMap["access_token"] = "create a token"
	dataMap["token_type"] = "Bearer or mac , baidu yi xia"
	// 原路返回
	dataMap["refresh_token"] = "create a refresh token"
	dataMap["expires_in"] = "expire"
	dataMap["scope"] = "rights"

	log.Println(dataMap)

	c.JSON(200, gin.H{
		"http": "true",
		"code": "0",
		"msg":  "",
		"data": dataMap,
	})
}
