package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"zqc/services"
)

func ParseInt(s string, min int, max int) (i int, err error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if i64 < int64(min) || i64 > int64(max) {
		return 0, errors.New(fmt.Sprintf("invalid value %v", i64))
	}
	return int(i64), nil
}

func ParseLocation(s string) (loc *services.Location, err error) {
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

	return &services.Location{float32(lon), float32(lat)}, nil
}

func ParseObjectId(s string) (id bson.ObjectId, err error) {
	if !bson.IsObjectIdHex(s) {
		return id, errors.New(fmt.Sprintf("invalid ObjectId %v", s))
	}
	return bson.ObjectIdHex(s), nil
}

func ParseObjectIds(s string) (ids []bson.ObjectId, err error) {
	ss := strings.Split(s, ",")
	for _, v := range ss {
		if !bson.IsObjectIdHex(v) {
			return nil, errors.New(fmt.Sprintf("invalid ObjectId %v", v))
		}
	}

	ids = make([]bson.ObjectId, 0, len(ss))
	for _, v := range ss {
		ids = append(ids, bson.ObjectIdHex(v))
	}
	return ids, nil
}

func ParseTime(s string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	return &t, err
}
