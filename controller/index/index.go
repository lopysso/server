package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeAction 默认页面
func HomeAction(c *gin.Context) {
	callbackUrl := c.DefaultQuery("callback", "/")
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"title":       "hello",
		"callbackUrl": callbackUrl,
	})
}
