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

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
}

type CreateUserRequest struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"passWord" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Permission int    `json:"permission" default:"1"`
}

type DeleteUserRequest struct {
	UserName string `json:"userName" binding:"required"`
}

type ChangePassWordRequest struct {
	UserName string `json:"userName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
}

type UpdateUserInfoRequest struct {
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Permission int    `json:"permission" binding:"required"`
}

type ListUserRequest struct{}

var requests = []Request{
	{
		"POST",
		"user.login",
		LoginRequest{},
		Login,
	},
	{
		"POST",
		"user.create",
		CreateUserRequest{},
		CreateUser,
	},
	{
		"POST",
		"user.delete",
		DeleteUserRequest{},
		DeleteUser,
	},
	{
		"POST",
		"user.update.password",
		ChangePassWordRequest{},
		ChangePassWord,
	},
	{
		"POST",
		"user.update.info",
		UpdateUserInfoRequest{},
		UpdateUserInfo,
	},
	{
		"GET",
		"user.list",
		ListUserRequest{},
		ListUser,
	},
}
