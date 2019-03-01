package handler

import (
	"reflect"
	. "webserver/args"
)

func (h *DefaultApiHandler) renderPermission(args *APIArgs, in interface{}) (out interface{}) {
	if reflect.ValueOf(in).Kind() == reflect.Slice {
		config, has := args.GetConfigTable(h.TableController.GetPermissionConfig())
		if has {
			access := h.apiHandlers.GetAccessConfig(args)
			if access == nil {
				return nil
			}
			out = config.InitTablePermission(in, access)
		} else {
			out = in
		}
	} else {
		out = in
	}
	return
}
