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
	m.Use(CORS)

	for _, v := range Neighborhoods {
		NeighborhoodMap[v.Id] = v.Name
	}

	m.Use(render.Renderer())
	m.Get("/", home)
	m.Get("/neighborhoods", getNeighborhoods)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}
}

func getNeighborhoods(w http.ResponseWriter, r render.Render) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := make(map[string]interface{})
	response["neighborhoods"] = Neighborhoods
	r.JSON(200, response)
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
