package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"time"
)

var Neighborhoods = map[int64]string{
	1:  "Chelsea",
	2:  "East Village / Lower East Side",
	3:  "Bushwick/Ridgewood",
	5:  "Tribeca / Downtown",
	6:  "SoHo",
	7:  "Williamsburg",
	8:  "Upper East Side",
	9:  "Midtown",
	10: "Gowanus",
	11: "Long Island City",
	12: "Museums",
	13: "Greenwich Village",
	14: "Dumbo",
	15: "Boerum Hill",
	16: "Bronx",
	17: "Sunset Park, Brooklyn",
	18: "Hell's Kitchen/Midtown West",
}

type Venue struct {
	Id             int64
	Name           string
	Address        string
	Website        string
	NeighborhoodId int64     `db:"neighborhood_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
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
}

type EventView struct {
	EventId             int64
	VenueId             int64
	Title               string
	VenueName           string
	VenueAddress        string
	VenueNeighborhoodId int64
	OpeningDate         time.Time
	OpeningStartTime    time.Time
	OpeningEndTime      time.Time
	Website             string
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	connection := os.Getenv("FILTERIZER_DSN")
	if connection == "" {
		connection = os.Getenv("DATABASE_URL")
	}
	db, err := sql.Open("postgres", connection)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Venue{}, "venues").SetKeys(true, "Id")
	dbmap.AddTableWithName(Event{}, "events").SetKeys(true, "Id")

	return dbmap
}

func OpeningSoon() []EventView {
	dbmap := initDb()
	var events []EventView
	const query = `
		select e.title Title, e.id EventId, v.id VenueId, v.name VenueName, v.address VenueAddress, 
		v.neighborhood_id VenueNeighborhoodId, e.opening_date OpeningDate, e.opening_start_time OpeningStartTime, 
		e.opening_end_time OpeningEndTime, coalesce(e.website, v.website) Website
		from events e, venues v 
		where e.venue_id = v.id and opening_date between current_date and current_date + interval '10 days' 
		order by opening_date, opening_start_time
	`
	_, err := dbmap.Select(&events, query)
	checkErr(err, "Event select failed")
	return events
}

func (e *EventView) OpeningDateTime() string {
	str := e.OpeningDate.Format("Monday, January 2, ")
	str += e.OpeningStartTime.Format("3:04")
	str += e.OpeningEndTime.Format("-3:04 PM")
	return strings.Replace(str, ":00", "", 2)
}

func (e *EventView) Neighborhood() string {
	return Neighborhoods[e.VenueNeighborhoodId]
}
