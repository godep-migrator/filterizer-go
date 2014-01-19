package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func CORS(c martini.Context, res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Methods", "*")
	c.Next()
}
