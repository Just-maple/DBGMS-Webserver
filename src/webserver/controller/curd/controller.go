package curd

import (
	. "webserver/args"
	"webserver/controller"
	"webserver/dbx"
	. "webserver/handler"
)

type CurdController struct {
	collection *dbx.CollectionController
	controller.DefaultController
}

func NewCurdController(c *dbx.Collection, in interface{}) *CurdController {
	return &CurdController{
		collection: c.CreateController(in),
	}
}

func (c *CurdController) Get(addr string, pm ...PermissionAuth) {
	c.RegisterGetApi(addr, func(args *APIArgs) (ret interface{}, err error) {
		return c.collection.GetAll(nil)
	}, pm...)
}

func (c *CurdController) Update(addr string, pm ...PermissionAuth) {
	c.RegisterPostApi(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.collection.UpdateByJson(args.Json)
	}, pm...)
}

func (c *CurdController) Delete(addr string, pm ...PermissionAuth) {
	c.RegisterPostApi(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.collection.RemoveByJson(args.Json)
	}, pm...)
}

func (c *CurdController) Default() interface{} {
	return c.collection.NewModel()
}

func (c *CurdController) New(addr string, pm ...PermissionAuth) {
	c.RegisterPostApi(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.collection.NewFromJson(args.Json)
	}, pm...)
}
