package main

import (
	"testing"
	"time"
)

func TestEventView(t *testing.T) {
	const yymmdd = "2006-01-02"
	const hhmm = "15:04"

	start, _ := time.Parse(yymmdd, "2014-01-03")
	end, _ := time.Parse(yymmdd, "2014-01-30")
	start_time, _ := time.Parse(hhmm, "18:00")
	end_time, _ := time.Parse(hhmm, "20:00")

	event := EventView{1, 2, "Big Show", "Gallery Hoggard Wagner", "601 West 26th Street", 1, start, start_time, end_time, end, ""}
	const expected = "Friday, January 3, 6-8 PM"
	if event.OpeningDateTime() != expected {
		t.Errorf("event.OpeningDateTime() = %v, want %v", event.OpeningDateTime(), expected)
	}
}
