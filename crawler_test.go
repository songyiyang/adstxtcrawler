package adstxtcrawler

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"testing"
)

func TestDo(t *testing.T) {

	logger := log.New(os.Stdout, "", 0)
	// could change this to random port instead of using the default
	mongo, err := InitDefaultMongoCxt()
	if err != nil {
		t.Error(err)
	}

	if err = mongo.DropDB(); err != nil {
		t.Error(err)
	}

	crawler := Crawler{nil, logger, mongo, &DefaultParser{logger}}

	for i, s := range []struct {
		publisherName string
		url           string
		expectedLen   int
	}{
		{"CNN", "http://www.cnn.com/ads.txt", 94},
		// there's one duplication: sharethrough.com, 9ef5c2ea, DIRECT, d53b998a7bd4ecd2
		{"Gizmodo", "http://www.gizmodo.com/ads.txt", 78},
		{"NYTimes", "http://www.nytimes.com/ads.txt", 15},
		{"Bloomberg", "https://www.bloomberg.com/ads.txt", 24},
		{"WordPress", "https://wordpress.com/ads.txt", 3642},
	} {
		if _, err = mongo.RecordsColl.RemoveAll(bson.M{}); err != nil {
			t.Errorf("t%d: %s", i, err.Error())
		}

		if err := crawler.Do(s.publisherName, s.url, nil); err != nil {
			t.Errorf("t%d: %s", i, err.Error())
		}

		count, err := mongo.RecordsColl.Find(bson.M{"publisherName": s.publisherName}).Count()
		if err != nil {
			t.Errorf("t%d: %s", i, err.Error())
		}

		if count != s.expectedLen {
			t.Errorf("t%d: %d wanted, got %d for publisher [%s]", i, s.expectedLen, count, s.publisherName)
		}
	}
}
