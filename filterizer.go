package main

import (
	"fmt"
	"log"
	// "github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	// "time"
)

func main() {
	http.HandleFunc("/", hello)
	fmt.Println("listening...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	dbmap := initDb()
	obj, err := dbmap.Get(Venue{}, 4)
	checkErr(err, "failed to load Venue")
	venue := obj.(*Venue)
	fmt.Fprintln(res, venue.Name)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
