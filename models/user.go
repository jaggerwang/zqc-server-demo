package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Mobile     string        `bson:"mobile"`
	Password   string        `bson:"password"`
	Salt       string        `bson:"salt"`
	Nickname   string        `bson:"nickname,omitempty"`
	Gender     string        `bson:"gender,omitempty"`
	CreateTime *time.Time    `bson:"createTime"`
	UpdateTime *time.Time    `bson:"updateTime,omitempty"`
}

type UserColl struct {
	*MongoColl
}

func NewUserColl() (uc *UserColl, err error) {
	coll, err := NewMongoColl("zqc", "zqc", "user")
	if err != nil {
		return nil, err
	}

	return &UserColl{coll}, nil
}
