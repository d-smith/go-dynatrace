package main

import (
	"net/http"

	"github.com/d-smith/go-dynatrace/sidecar/server"
)

//Serve up some http goodness
func serveHttp() {
	handleCallHandler := http.HandlerFunc(server.HandleCall)
	http.ListenAndServe(":8080", handleCallHandler)
}

func main() {
	serveHttp()
}
