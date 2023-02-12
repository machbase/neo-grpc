package spi

import "strings"

// Refer: https://gosamples.dev/date-time-format-cheatsheet/
var _timeformats = map[string]string{
	"-":           "2006-01-02 15:04:05.999",
	"DEFAULT":     "2006-01-02 15:04:05.999",
	"NUMERIC":     "01/02 03:04:05PM '06 -0700", // The reference time, in numerical order.
	"ANSIC":       "Mon Jan _2 15:04:05 2006",
	"UNIX":        "Mon Jan _2 15:04:05 MST 2006",
	"RUBY":        "Mon Jan 02 15:04:05 -0700 2006",
	"RFC822":      "02 Jan 06 15:04 MST",
	"RFC822Z":     "02 Jan 06 15:04 -0700", // RFC822 with numeric zone
	"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
	"RFC1123":     "Mon, 02 Jan 2006 15:04:05 MST",
	"RFC1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with numeric zone
	"RFC3339":     "2006-01-02T15:04:05Z07:00",
	"RFC3339NANO": "2006-01-02T15:04:05.999999999Z07:00",
	"KITCHEN":     "3:04:05PM",
	"STAMP":       "Jan _2 15:04:05",
	"STAMPMILLI":  "Jan _2 15:04:05.000",
	"STAMPMICRO":  "Jan _2 15:04:05.000000",
	"STAMPNANO":   "Jan _2 15:04:05.000000000",
}

func GetTimeformat(f string) string {
	if m, ok := _timeformats[strings.ToUpper(f)]; ok {
		return m
	}
	return f
}
