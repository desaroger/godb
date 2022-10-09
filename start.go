package main

import "godb/http"

func main() {
	api := http.NewHttpJsonApi("_data")
	api.Start("localhost:5001")
}
