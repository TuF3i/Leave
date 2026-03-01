package geoIP

import (
	"sync"

	"github.com/ip2location/ip2location-go/v9"
)

type Location struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TimeZone    string  `json:"time_zone"`
}

type GeoIP struct {
	db *ip2location.DB
	mu sync.RWMutex
}

func GetGeoIP() (*GeoIP, error) {
	db, err := ip2location.OpenDB("./geoIP/IP2LOCATION-LITE-DB11.BIN")
	if err != nil {
		return nil, err
	}

	return &GeoIP{db: db, mu: sync.RWMutex{}}, nil
}
