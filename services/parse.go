package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func ParseInt(s string, min interface{}, max interface{}) (int, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("invalid int %v", s))
	}
	i := int(i64)
	if min != nil && i < min.(int) {
		return 0, errors.New(fmt.Sprintf("must >= %v", min))
	}
	if max != nil && i > max.(int) {
		return 0, errors.New(fmt.Sprintf("must <= %v", max))
	}
	return i, nil
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
