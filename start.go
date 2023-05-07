package main

import (
	"godb/http"
	logs "godb/logs"
)

func main() {
	logs.Initialize()

	addr := "localhost:5001"
	api := http.NewHttpJsonApi("_data")
	logs.Info("HttpJsonApi listening at %s\n", addr)

	api.Start(addr)
}
