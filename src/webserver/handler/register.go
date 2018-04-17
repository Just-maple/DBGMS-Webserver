package handler

type RegisterGroup struct {
	Route *JsonAPIFuncRoute
	pm    *[]PermissionAuth
}

func (r *RegisterGroup) RegisterDefaultAPI(name string, function DefaultAPIFunc) {
	r.Route.RegisterDefaultAPI(name, function, *r.pm...)
}
func (r *RegisterGroup) RegisterAPI(name string, function JsonAPIFunc) {
	r.Route.RegisterAPI(name, function, *r.pm...)
}

func (j *JsonAPIFuncRoute) MakeRegisterGroup(pm ...PermissionAuth) *RegisterGroup {
	return &RegisterGroup{j, &pm}
}

func (j *JsonAPIFuncRoute) RegisterAPI(name string, function JsonAPIFunc, pm ...PermissionAuth) {
	j.registerJsonAPI(name,
		func(args *APIArgs) (i interface{}, e error) {
			return function(args.context, args.json, args.session)
		}, pm)
}

func (j JsonAPIFuncRoute) registerJsonAPI(name string, function DefaultAPIFunc, pm []PermissionAuth) {
	if j[name] != nil {
		panic("route already existed")
	} else {
		j[name] = &DefaultAPI{
			function, pm,
		}
	}
}

func (j *JsonAPIFuncRoute) RegisterDefaultAPI(name string, api DefaultAPIFunc, pm ...PermissionAuth) {
	j.registerJsonAPI(name, api, pm)
}
