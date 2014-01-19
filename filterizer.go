package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	"os"
)

var NeighborhoodMap = make(map[int64]string)

func main() {
	m := martini.Classic()
	dbmap := initDb()
	m.Map(dbmap)

	for _, v := range Neighborhoods {
		NeighborhoodMap[v.Id] = v.Name
	}

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

func home(dbmap *gorp.DbMap, r render.Render) {
	tmpl_vars := make(map[string]interface{})
	tmpl_vars["openingSoon"] = openingSoon(dbmap)
	tmpl_vars["openNow"] = openNow(dbmap)
	r.HTML(200, "home", tmpl_vars)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
