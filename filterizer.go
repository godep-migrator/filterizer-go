package main

import (
	// "fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"log"
	"net/http"
	"os"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/", home)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}

func home(r render.Render) {
	// var tmpl_vars map[string]interface{}
	tmpl_vars := make(map[string]interface{})
	tmpl_vars["opening_soon"] = OpeningSoon()
	r.HTML(200, "home", tmpl_vars)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
