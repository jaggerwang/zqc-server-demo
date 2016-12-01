package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type File struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Service    string        `bson:"service"`
	Region     string        `bson:"region"`
	Bucket     string        `bson:"bucket"`
	ObjectKey  string        `bson:"objectKey"`
	Size       int           `bson:"size"`
	Filename   string        `bson:"filename,omitempty"`
	UploaderId bson.ObjectId `bson:"uploaderId,omitempty"`
	CreateTime *time.Time    `bson:"createTime"`
	UpdateTime *time.Time    `bson:"updateTime,omitempty"`
}

type FileColl struct {
	*MongoColl
}

func NewFileColl() (fc *FileColl, err error) {
	coll, err := NewMongoColl("zqc", "zqc", "file")
	if err != nil {
		return nil, err
	}

	return &FileColl{coll}, nil
}
