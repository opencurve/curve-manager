package user

import (
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	MODULE = "user"
)

var METHOD_REQUEST map[string]Request

type (
	HandlerFunc func(r *pigeon.Request, ctx *Context) bool

	Context struct {
		Data interface{}
	}

	Request struct {
		httpMethod string
		method     string
		vType      interface{}
		handler    HandlerFunc
	}
)

func init() {
	METHOD_REQUEST = map[string]Request{}
	for _, request := range requests {
		METHOD_REQUEST[request.method] = request
		core.BELONG[request.method] = MODULE
	}
}

type CreateUserRequest struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"passWord" binding:"required"`
	Permission int    `json:"permission" binding:"required"`
}

var requests = []Request{
	{
		"POST",
		"user.add",
		CreateUserRequest{},
		CreateUser,
	},
}
