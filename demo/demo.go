package main

import (
	"webserver"
)

func main() {
	//get new web-server container from your handler and your config
	svr := webserver.NewWebServerFromHandler(NewConfig(), new(ApiHandler))
	svr.Start()
}
