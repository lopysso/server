package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lopysso/server/controller"
	"github.com/lopysso/server/libs/oauth/token"
)

func Middleware(c *gin.Context)  {
	// token
	query := queryParams{}
	jsonRes := controller.NewJsonRes()

	err := c.ShouldBindQuery(&query)

	if err != nil{
		c.JSON(http.StatusOK,jsonRes.Error(query.getError(err.(validator.ValidationErrors))).Generate())
		c.Abort()
		return
	}

	// get from db
	acc,err := token.GetAccessAvailableFromDb(query.AccessToken)
	if err != nil{
		c.JSON(http.StatusOK,jsonRes.Error(err.Error()).Generate())
		c.Abort()
		return
	}

	c.Set("accessToken",acc)
	c.Next()
}

type queryParams struct {
	AccessToken string `form:"access_token" binding:"required,alphanum,min=40,max=40"`
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


// 后面再说
var base58Validator validator.Func = func(fl validator.FieldLevel) bool {
	
	log.Println("custom base58 ： ",fl.Field())
	
	return true
}