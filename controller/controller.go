package controller

import "github.com/gin-gonic/gin"

type JsonRes struct {
	http bool
	code int
	msg  string
	data gin.H
}

func NewJsonRes() JsonRes {
	j := JsonRes{}
	j.http = true
	j.code = 0
	j.msg = ""
	j.data = gin.H{}
	return j
}

func (p *JsonRes) Data(data gin.H) *JsonRes {
	p.data = data
	return p
}

func (p *JsonRes) Success(msg ...string) *JsonRes {
	if len(msg) > 0 {
		p.msg = msg[0]
	} else {
		p.msg = "ok"
	}
	p.code = 0
	return p
}

func (p *JsonRes) Error(msg string, code ...int) *JsonRes {
	if len(code) > 0 {
		p.code = code[0]
	} else {
		p.code = 1
	}
	p.msg = msg
	return p
}

func (p *JsonRes) Generate() gin.H {
	return gin.H{
		"http": p.http,
		"code": p.code,
		"msg":  p.msg,
		"data": p.data,
	}
}
