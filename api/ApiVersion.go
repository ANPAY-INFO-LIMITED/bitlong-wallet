package api

import (
	"time"
)

const (
	PrimaryVersion = "0"
	SubVersion     = "2"
	DevVersion     = "93"
)

const (
	BaseVersion = "v" + PrimaryVersion + "." + SubVersion + "." + DevVersion
	DateTime    = "2025-10-16 17:30:00"
)

func apiVersionWithMaker(maker string) string {
	dt := DateTime

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-1-02 15:04:05",
		"2006-01-2 15:04:05",
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, dt); err == nil {
			dt = t.Format("20060102150405")
		}
	}

	if maker == "" {
		return BaseVersion + "_" + dt
	} else {
		return BaseVersion + "_" + maker + "_" + dt
	}
}

func GetApiVersion() string {
	return apiVersionWithMaker("")
}
