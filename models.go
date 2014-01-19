package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"time"
)

type Neighborhood struct {
	Id   int64
	Name string
}

var Neighborhoods = []Neighborhood{
	{15, "Boerum Hill"},
	{16, "Bronx"},
	{3, "Bushwick/Ridgewood"},
	{1, "Chelsea"},
	{14, "Dumbo"},
	{2, "East Village / Lower East Side"},
	{10, "Gowanus"},
	{13, "Greenwich Village"},
	{18, "Hell's Kitchen/Midtown West"},
	{11, "Long Island City"},
	{9, "Midtown"},
	{12, "Museums"},
	{19, "Park Slope"},
	{6, "SoHo"},
	{17, "Sunset Park, Brooklyn"},
	{5, "Tribeca / Downtown"},
	{8, "Upper East Side"},
	{7, "Williamsburg"},
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
	EndDate             time.Time
	Website             string
}

// struct used for showing events organized by neighborhood
type NeighborhoodEvents struct {
	Neighborhood string
	Events       []EventView
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

const eventViewFields = `
	select e.title Title, e.id EventId, v.id VenueId, v.name VenueName, v.address VenueAddress, 
	v.neighborhood_id VenueNeighborhoodId, coalesce(e.opening_date, 'epoch'::date) OpeningDate, e.opening_start_time OpeningStartTime, 
	e.opening_end_time OpeningEndTime, e.end_date EndDate, coalesce(e.website, v.website) Website
	from events e, venues v
`

func openingSoon(dbmap *gorp.DbMap) []EventView {
	var events []EventView
	const query = eventViewFields + `
		where e.venue_id = v.id and opening_date between date(current_date at time zone 'EST') and date(current_date at time zone 'EST') + interval '10 days' 
		order by opening_date, opening_start_time
	`
	_, err := dbmap.Select(&events, query)
	checkErr(err, "openingSoon select failed")
	return events
}

func openNow(dbmap *gorp.DbMap) []NeighborhoodEvents {
	list := make([]NeighborhoodEvents, 0, len(Neighborhoods))
	for _, value := range Neighborhoods {
		events := openByNeighborhood(dbmap, value.Id)
		if len(events) > 0 {
			list = append(list, NeighborhoodEvents{value.Name, events})
		}
	}
	return list
}

func openByNeighborhood(dbmap *gorp.DbMap, hood_id int64) []EventView {
	var events []EventView
	const query = eventViewFields + `
		where e.venue_id = v.id and e.start_date <= date(current_date at time zone 'EST')
		  and e.end_date >= date(current_date at time zone 'EST') and v.neighborhood_id = $1
		order by e.end_date
	`
	_, err := dbmap.Select(&events, query, hood_id)
	checkErr(err, "openByNeighborhood select failed")
	return events
}

func (e *EventView) OpeningDateTime() string {
	str := e.OpeningDate.Format("Monday, January 2, ")
	str += e.OpeningStartTime.Format("3:04")
	str += e.OpeningEndTime.Format("-3:04 PM")
	return strings.Replace(str, ":00", "", 2)
}

func (e *EventView) FormattedEndDate() string {
	return e.EndDate.Format("Monday, January 2")
}

func (e *EventView) Neighborhood() string {
	return NeighborhoodMap[e.VenueNeighborhoodId]
}
