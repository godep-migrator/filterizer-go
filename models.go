package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/lib/pq"
	"os"
	"strings"
	"time"
)

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	connection := os.Getenv("FILTERIZER_DSN")
	if connection == "" {
		url := os.Getenv("DATABASE_URL")
		connection, _ := pq.ParseURL(url)
		connection += " sslmode=require"
	}
	db, err := sql.Open("postgres", connection)
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
	Id                 int64
	Name               string
	Address            string
	Website            string
	NeighborhoodId     int64         `db:"neighborhood_id"`
	CreatedAt          time.Time     `db:"created_at"`
	UpdatedAt          time.Time     `db:"updated_at"`
	CachedNeighborhood *Neighborhood `db:"-"`
}

type Event struct {
	Id               int64
	Title            string
	VenueId          int64     `db:"venue_id"`
	StartDate        time.Time `db:"start_date"`
	EndDate          time.Time `db:"end_date"`
	OpeningDate      time.Time `db:"opening_date"`
	OpeningStartTime time.Time `db:"opening_start_time"`
	OpeningEndTime   time.Time `db:"opening_end_time"`
	Website          string
	Tweeted          bool
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	CachedVenue      *Venue    `db:"-"`
}

func OpeningSoon() []Event {
	dbmap := initDb()
	var events []Event
	_, err := dbmap.Select(&events,
		"select * from events where opening_date between current_date and current_date + interval '10 days' order by opening_date, opening_start_time")
	checkErr(err, "Event select failed")
	return events
}

func (e *Event) OpeningDateTime() string {
	str := e.OpeningDate.Format("Monday, January 2, ")
	str += e.OpeningStartTime.Format("3:04")
	str += e.OpeningEndTime.Format("-3:04 PM")
	return strings.Replace(str, ":00", "", 2)
}

func (e *Event) Url() string {
	if e.Website != "" {
		return e.Website
	}
	return e.Venue().Website
}

func (e *Event) Venue() *Venue {
	if e.CachedVenue != nil {
		return e.CachedVenue
	}
	dbmap := initDb()
	obj, err := dbmap.Get(Venue{}, e.VenueId)
	checkErr(err, "Venue select failed")
	venue := obj.(*Venue)
	e.CachedVenue = venue
	return venue
}

func (v *Venue) Neighborhood() *Neighborhood {
	if v.CachedNeighborhood != nil {
		return v.CachedNeighborhood
	}
	dbmap := initDb()
	obj, err := dbmap.Get(Neighborhood{}, v.NeighborhoodId)
	checkErr(err, "Neighborhood select failed")
	neighborhood := obj.(*Neighborhood)
	v.CachedNeighborhood = neighborhood
	return neighborhood
}
