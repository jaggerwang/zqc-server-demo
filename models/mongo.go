package models

import (
	"math"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

var MongoSessions = map[string]*mgo.Session{}

func NewMongoSession(clusterName string) (session *mgo.Session, err error) {
	if session, ok := MongoSessions[clusterName]; !ok {
		config := viper.GetStringMap("mongodb." + clusterName)
		info := mgo.DialInfo{
			Addrs: strings.Split(config["addrs"].(string), ","),
		}
		timeout, ok := config["timeout"].(int64)
		if ok {
			info.Timeout = time.Duration(timeout * int64(math.Pow10(9)))
		}
		session, err = mgo.DialWithInfo(&info)
		if err != nil {
			return nil, err
		}
		session.SetMode(mgo.Monotonic, false)
		session.SetSafe(&mgo.Safe{
			WMode: "majority",
		})
		MongoSessions[clusterName] = session
	}
	return MongoSessions[clusterName].Copy(), nil
}

type MongoDB struct {
	*mgo.Database
}

func NewMongoDB(clusterName string, dbName string) (db *MongoDB, err error) {
	session, err := NewMongoSession(clusterName)
	if err != nil {
		return nil, err
	}
	return &MongoDB{session.DB(dbName)}, nil
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

type MongoColl struct {
	*mgo.Collection
}

func NewMongoColl(clusterName string, dbName string, collName string) (coll *MongoColl, err error) {
	db, err := NewMongoDB(clusterName, dbName)
	if err != nil {
		return nil, err
	}
	return &MongoColl{db.C(collName)}, nil
}

func (m *MongoColl) Close() {
	m.Database.Session.Close()
}

func EmptyDB(clusterName string, dbName string, collName string) (err error) {
	var collNames []string
	if collName == "" {
		collNames, err = DBCollNames(clusterName, dbName)
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

		_, err = coll.RemoveAll(nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func DBCollNames(clusterName string, dbName string) (collNames []string, err error) {
	db, err := NewMongoDB(clusterName, dbName)
	if err != nil {
		return nil, err
	}
	collNames, err = db.CollectionNames()
	if err != nil {
		return nil, err
	}
	return collNames, nil
}
