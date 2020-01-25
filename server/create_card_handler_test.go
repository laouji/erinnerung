package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/laouji/erinnerung/server"
)

func TestEstimateDueDate(t *testing.T) {
	testCases := map[string]string{
		"Mon, 03 Feb 2020 15:04:05 UTC": "Mon, 10 Feb 2020 11:00:00 UTC",
		"Tue, 07 Jan 2020 18:18:18 UTC": "Mon, 13 Jan 2020 11:00:00 UTC",
		"Wed, 01 Jan 2020 18:18:18 UTC": "Mon, 06 Jan 2020 11:00:00 UTC",
		"Thu, 02 Jan 2020 18:18:18 UTC": "Mon, 06 Jan 2020 11:00:00 UTC",
		"Fri, 03 Jan 2020 00:00:00 UTC": "Mon, 06 Jan 2020 11:00:00 UTC",
		"Sat, 30 May 2020 23:59:59 UTC": "Mon, 01 Jun 2020 11:00:00 UTC",
		"Sun, 27 Sep 2020 07:07:07 UTC": "Mon, 28 Sep 2020 11:00:00 UTC",
	}

	for todayStr, expectedDateStr := range testCases {
		today, err := time.Parse(time.RFC1123, todayStr)
		if err != nil {
			fmt.Printf("could not parse time: %s\n", err)
			t.Fail()
		}

		nextMonday := server.EstimateDueDate(today)

		nextMondayStr := nextMonday.Format(time.RFC1123)
		if nextMondayStr != expectedDateStr {
			fmt.Printf("expected next monday to be %s but was %s\n", expectedDateStr, nextMondayStr)
			t.Fail()
		}

		if nextMonday.Weekday() != 1 {
			fmt.Printf("resulting next monday was not a monday: %s\n", nextMonday)
			t.Fail()
		}
	}
}
