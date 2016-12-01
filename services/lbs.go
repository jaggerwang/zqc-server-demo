package services

import (
	"errors"
	"strconv"
	"strings"
)

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

func ParseLocation(s string) (loc *Location, err error) {
	slice := strings.Split(s, ",")
	if len(slice) != 2 {
		return nil, errors.New("invalid location")
	}

	lon, err := strconv.ParseFloat(slice[0], 32)
	if err != nil {
		return nil, errors.New("invalid longitude")
	}
	lat, err := strconv.ParseFloat(slice[1], 32)
	if err != nil {
		return nil, errors.New("invalid latitude")
	}

	return &Location{float32(lon), float32(lat)}, nil
}
