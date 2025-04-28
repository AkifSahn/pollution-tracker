package pollution

import (
	"strconv"
	"time"
)

var TimeFormat string = "2006-01-02 15:04:05"

func ParseTimeRange(fromStr, toStr string, from, to *time.Time) (bool, string) {
	var err error
	*from, err = time.Parse(TimeFormat, fromStr)
	if err != nil {
		return false, "Incorrect time format!"
	}

	*to, err = time.Parse(TimeFormat, toStr)
	if err != nil {
		return false, "Incorrect time format!"
	}

	return true, ""
}

func ParseLatLon(latStr, lonStr string, lat, lon *float64) (bool, string) {
	var err error
	*lat, err = strconv.ParseFloat(latStr, 64)
	if err != nil {
		return false, "Incorrect latitude format!"
	}

	*lon, err = strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return false, "Incorrect longitude format!"
	}

	return true, ""
}
