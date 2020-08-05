package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller/index"
	"github.com/lopysso/server/controller/oauth"
	oauthApi "github.com/lopysso/server/controller/oauth/api"
	"github.com/lopysso/server/controller/user"
)

// NewRouter 生成默认路由
func NewRouter() *gin.Engine {

	routerDefault := gin.Default()

	// r.StaticFile("/index.html", "../public/index.html")
	routerDefault.Static("/assets", "../public/assets")

	// r.LoadHTMLFiles("view/*")
	// 这个位置要注意
	routerDefault.LoadHTMLGlob("../view/**/*")

	routerDefault.Use(func(c *gin.Context) {
		c.Header("Pragma", "no-cache")
	})
	// login
	routerDefault.GET("/", index.HomeAction)
	routerDefault.POST("/signin", index.SigninAction)

	//
	// routerDefault.Group("/oauth", func(c *gin.Context){

	// })

	//

	routerDefault.GET("/oauth/authorize", user.Middleware, user.Authorize)

	//
	routerDefault.GET("/oauth/access_token", oauth.AccessToken)
	routerDefault.GET("/oauth/refresh_token", oauth.RefreshTokenAction)

	// oauth api
	oauthApiRouter := routerDefault.Group("/oauth/api")
	oauthApiRouter.Use(oauthApi.Middleware)
	{
		oauthApiRouter.GET("/userinfo", oauthApi.Userinfo)
	}

	//
	userRouter := routerDefault.Group("/user")
	userRouter.Use(user.Middleware)
	{
		userRouter.GET("/", user.HomeAction)
	}
	// userRouter :=routerDefault.GET("/user/", user.Middleware, user.HomeAction)

	return routerDefault
}
