package handler

import "webserver/dbx"

type CurdController struct {
	controller *dbx.CollectionController
	getRoute   *JsonAPIFuncRoute
	postRoute  *JsonAPIFuncRoute
}

func (h *DefaultApiHandler) RegisterCurdAPI(c *dbx.Collection, in interface{}) (*CurdController) {
	return &CurdController{
		c.CreateController(in),
		&h.ApiGetHandlers,
		&h.ApiPostHandlers,
	}
}

func (c *CurdController) Get(addr string, pm ...PermissionAuth) {
	c.getRoute.RegisterDefaultAPI(addr, func(args *APIArgs) (ret interface{}, err error) {
		return c.controller.GetAll(nil)
	}, pm...)
}

func (c *CurdController) Update(addr string, pm ...PermissionAuth) {
	c.postRoute.RegisterDefaultAPI(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.controller.UpdateByJson(args.json)
	}, pm...)
}

func (c *CurdController) Delete(addr string, pm ...PermissionAuth) {
	c.postRoute.RegisterDefaultAPI(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.controller.RemoveByJson(args.json)
	}, pm...)
}

func (c *CurdController) Default() interface{} {
	return c.controller.NewModel()
}

func (c *CurdController) New(addr string, pm ...PermissionAuth) {
	c.postRoute.RegisterDefaultAPI(addr, func(args *APIArgs) (ret interface{}, err error) {
		return nil, c.controller.NewFromJson(args.json)
	}, pm...)
}
