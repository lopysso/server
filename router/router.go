package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lopysso/server/controller/index"
	"github.com/lopysso/server/controller/oauth"
	"github.com/lopysso/server/libs/account"
)

// NewRouter 生成默认路由
func NewRouter() *gin.Engine {

	routerDefault := gin.Default()

	// r.StaticFile("/index.html", "../public/index.html")
	routerDefault.Static("/assets", "../public/assets")

	// r.LoadHTMLFiles("view/*")
	// 这个位置要注意
	routerDefault.LoadHTMLGlob("../view/**/*")

	// login
	routerDefault.GET("/", index.HomeAction)
	routerDefault.POST("/signin", index.SigninAction)

	//
	// routerDefault.Group("/oauth", func(c *gin.Context){

	// })

	oauthRouter := routerDefault.Group("/oauth", func(c *gin.Context) {
		log.Println("hello oauth start")
		token, err := c.Cookie("authorize_token")
		if err != nil {
			c.Redirect(302, "/")
			log.Println("wheratasfasdfaslkdfjasdfasodfjoasdjfasdf")
		}

		// token
		mo,err := account.GetFromToken(token)
		if err != nil {
			c.Redirect(302, "/")
			log.Println("wheratasfasdfaslkdfjasdfasodfjoasdjfasdf")
		}

		c.Set("userInfo", mo)



		// find user info

		c.Next()
		log.Println("hello oauth end")
	})
	oauthRouter.GET("/authorize", oauth.Authorize)

	//
	routerDefault.GET("/oauth2/access_token", oauth.AccessToken)

	return routerDefault
}
