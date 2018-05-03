package handler

import (
	"reflect"
	. "webserver/args"
	"webserver/permission"
)

func (h *DefaultApiHandler) GetAccessConfigFromArgs(arg *APIArgs) (access permission.AccessConfig) {
	return h.GetAccessConfig(arg)
}

func (h *DefaultApiHandler) renderPermission(args *APIArgs, in interface{}) (out interface{}) {
	if reflect.ValueOf(in).Kind() == reflect.Slice {
		config, has := args.GetConfigTable(h.TableController.GetPermissionConfig())
		if has {
			access := h.GetAccessConfigFromArgs(args)
			out = config.InitTablePermission(in, access)
		} else {
			out = in
		}
	} else {
		out = in
	}
	return
}
