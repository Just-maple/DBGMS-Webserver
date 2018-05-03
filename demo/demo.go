package main

import (
	"webserver"
)

func main() {
	//init custom handler
	h := new(ApiHandler)
	//new handler by config and handler
	svr := webserver.NewWebServerFromHandler(NewConfig(), h)
	svr.Start()
}
