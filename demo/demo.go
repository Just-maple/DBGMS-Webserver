package main

import (
	"webserver"
)

func main() {
	svr := webserver.NewWebServerFromHandler(NewConfig(), new(ApiHandler))
	svr.Start()
}
