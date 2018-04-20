package main

import (
	"webserver"
	"syncx"
)

func main() {
	syncx.Test()
	h := new(ApiHandler)
	svr := webserver.NewWebServerFromHandler(NewConfig(), h)
	svr.Start()
}
