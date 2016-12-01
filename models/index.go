package models

import (
	"gopkg.in/mgo.v2"
)

var ZqcDbIndexes = map[string][]mgo.Index{
	"user": []mgo.Index{
		mgo.Index{
			Key:        []string{"username"},
			Unique:     true,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"nickname"},
			Unique:     true,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"mobile"},
			Unique:     true,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"email"},
			Unique:     true,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"$2dsphere:location"},
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"file": []mgo.Index{
		mgo.Index{
			Key:        []string{"uploaderId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"court": []mgo.Index{
		mgo.Index{
			Key:        []string{"$2dsphere:location"},
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"post": []mgo.Index{
		mgo.Index{
			Key:        []string{"creatorId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"courtId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"$2dsphere:location"},
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"cityCode", "sportCode"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"post_like": []mgo.Index{
		mgo.Index{
			Key:        []string{"postId", "userId"},
			Unique:     true,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"userId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"post_comment": []mgo.Index{
		mgo.Index{
			Key:        []string{"postId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"userId"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"createtime"},
			Unique:     false,
			Sparse:     false,
			Background: true,
		},
	},
	"user_stat": []mgo.Index{
		mgo.Index{
			Key:        []string{"post"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"likePost"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"commentPost"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
	},
	"post_stat": []mgo.Index{
		mgo.Index{
			Key:        []string{"liked"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"commented"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
	},
	"court_stat": []mgo.Index{
		mgo.Index{
			Key:        []string{"user"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
		mgo.Index{
			Key:        []string{"post"},
			Unique:     false,
			Sparse:     true,
			Background: true,
		},
	},
}

func CreateDbIndexes(clusterName string, dbName string, collName string, pos int) (err error) {
	var collNames []string
	if collName == "" {
		collNames, err = DbCollNames(clusterName, dbName)
		if err != nil {
			return err
		}
	} else {
		collNames = []string{collName}
	}

	for _, collName := range collNames {
		coll, err := NewMongoColl(clusterName, dbName, collName)
		if err != nil {
			return err
		}

		for i, index := range ZqcDbIndexes[collName] {
			if pos == -1 || i == pos {
				err := coll.EnsureIndex(index)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
