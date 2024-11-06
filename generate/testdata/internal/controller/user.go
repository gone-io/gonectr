package controller

import (
	"test/internal/interface/entity"
	"test/internal/interface/service"

	"github.com/gone-io/gone"
)

type userCtr struct {
	gone.Flag

	root  gone.RouteGroup `gone:"*"`
	iCook service.ICook   `gone:"*"`
	iUser service.IUser   `gone:"*"`
}

func (ctr *userCtr) Mount() gone.GinMountError {
	ctr.root.
		POST("/users", func(in struct {
			req entity.User `gone:"http,body"`
		}) error {
			return ctr.iUser.CreateUser(&in.req)
		}).
		GET("/users", ctr.iUser.GetUserList).
		POST("/cooks", func(in struct {
			req entity.Cook `gone:"http,body"`
		}) error {
			return ctr.iCook.CreateCook(&in.req)
		}).
		GET("/cooks", ctr.iCook.GetCookList)
	return nil
}
