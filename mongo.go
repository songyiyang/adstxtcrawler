package adstxtcrawler

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	defaultMongoHost = "127.0.0.1"
	defaultMongoPort = "27017"

	dbName                = "adstxtcrawler"
	recordsCollectionName = "records"
)

type MongoCxt struct {
	Mgo *mgo.Session
	// db.getCollection("records").createIndex({ "publisherName": 1, "domainName": 1, "publisherAccountId": 1, "accountType": 1})
	RecordsColl *mgo.Collection
}

func InitDefaultMongoCxt() (*MongoCxt, error) {

	maxWait := time.Duration(10 * time.Second)
	session, err := mgo.DialWithTimeout(fmt.Sprintf("%s:%s", defaultMongoHost, defaultMongoPort), maxWait)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Strong, true)

	return &MongoCxt{session, session.DB(dbName).C(recordsCollectionName)}, nil
}

// for multiple dbs, could parse dbname as args
func (s *MongoCxt) DropDB() error {
	return s.Mgo.DB(dbName).DropDatabase()
}
