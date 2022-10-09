package main

import "godb/http"

func main() {
	api := http.NewHttpJsonApi()
	api.Start("localhost:5001")
}
