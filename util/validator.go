package util

import (
	valid "github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	valid.TagMap["objectidhex"] = bson.IsObjectIdHex
}
