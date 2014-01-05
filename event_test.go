package main

import (
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	const yymmdd = "2006-01-02"
	const hhmm = "15:04"

	start, _ := time.Parse(yymmdd, "2014-01-03")
	end, _ := time.Parse(yymmdd, "2014-01-30")
	start_time, _ := time.Parse(hhmm, "18:00")
	end_time, _ := time.Parse(hhmm, "20:00")

	event := Event{1, "Big Show", 2, start, end, start, start_time, end_time, "", false, start, start, nil}
	const expected = "Friday, January 3, 6-8 PM"
	if event.OpeningDateTime() != expected {
		t.Errorf("event.OpeningDateTime() = %v, want %v", event.OpeningDateTime(), expected)
	}
}
