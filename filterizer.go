package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"log"
	"net/http"
	"os"
)

var NeighborhoodMap = make(map[int64]string)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/", home)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	for _, v := range Neighborhoods {
		NeighborhoodMap[v.Id] = v.Name
	}

	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}

func home(r render.Render) {
	tmpl_vars := make(map[string]interface{})
	dbmap := initDb()
	tmpl_vars["openingSoon"] = openingSoon(dbmap)
	tmpl_vars["openNow"] = openNow(dbmap)
	r.HTML(200, "home", tmpl_vars)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
