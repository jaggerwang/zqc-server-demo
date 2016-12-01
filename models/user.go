package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	Username       string        `bson:"username"`
	Password       string        `bson:"password"`
	Salt           string        `bson:"salt"`
	Nickname       string        `bson:"nickname,omitempty"`
	Gender         string        `bson:"gender,omitempty"`
	Mobile         string        `bson:"mobile,omitempty"`
	AvatarType     string        `bson:"avatarType"`
	AvatarName     string        `bson:"avatarName,omitempty"`
	AvatarId       bson.ObjectId `bson:"avatarId,omitempty"`
	Email          string        `bson:"email,omitempty"`
	Intro          string        `bson:"intro,omitempty"`
	BackgroundType string        `bson:"backgroundType"`
	BackgroundName string        `bson:"backgroundName,omitempty"`
	BackgroundId   bson.ObjectId `bson:"backgroundId,omitempty"`
	Location       *Point        `bson:"location"`
	CreateTime     *time.Time    `bson:"createTime,omitempty"`
	UpdateTime     *time.Time    `bson:"updateTime,omitempty"`
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
