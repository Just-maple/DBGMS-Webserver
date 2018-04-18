package main

import (
	"webserver"
)

func main() {
	h := new(ApiHandler)
	svr := webserver.NewWebServerFromHandler(NewConfig(), h)
	svr.Start()
}
