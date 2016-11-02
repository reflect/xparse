package xtime

import (
	"time"
)

var (
	base   = "2016-01-01 05:00:00"
	format = "2006-01-02 15:04:05"
)

func TimeZoneOffset(offset string) (hours int, minutes int, err error) {
	loc, err := time.LoadLocation(offset)

	if err != nil {
		return 0, 0, err
	}

	utc, _ := time.ParseInLocation(format, base, time.UTC)
	local, _ := time.ParseInLocation(format, base, loc)

	diff := utc.Sub(local)

	hours = int(diff / time.Hour)
	minutes = int((diff - (time.Duration(hours) * time.Hour)) / time.Minute)

	return
}
