package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", os.Getenv("FILTERIZER_DSN"))
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Neighborhood{}, "neighborhoods").SetKeys(true, "Id")
	dbmap.AddTableWithName(Venue{}, "venues").SetKeys(true, "Id")
	dbmap.AddTableWithName(Event{}, "events").SetKeys(true, "Id")

	return dbmap
}

type Neighborhood struct {
	Id   int64
	Name string
}

type Venue struct {
	Id            int64
	Name          string
	Address       string
	Website       string
	NeigborhoodId int64 `db:"neighborhood_id"`
}

type Event struct {
	Id               int64
	Title            string
	VenueId          int64 ` db:"venue_id"`
	StartDate        time.Time
	EndDate          time.Time
	OpeningDate      time.Time
	OpeningStartTime time.Time
	OpeningEndTime   time.Time
	Website          string
}
