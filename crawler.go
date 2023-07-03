package adstxtcrawler

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

const bulkMaxWriteSizePerBatch = 1000

type Crawler struct {
	WaitGroup *sync.WaitGroup
	Log       *log.Logger
	Mongo     *MongoCxt
	Parser    Parser
}

// Do will take a url, it's header and publisher name, retrieve the publisher's ads.txt list, parse the list and persist
func (s *Crawler) Do(publisherName, url string, headers map[string]string) error {
	s.Log.Printf("start crawling publisher[%s] with given url[%s]...", publisherName, url)
	defer func() {
		s.Log.Printf("finish crawling publisher[%s] with given url[%s]", publisherName, url)

		if s.WaitGroup != nil {
			s.WaitGroup.Done()
		}
	}()

	body, err := s.getURLBody(url, headers)
	if err != nil {
		panic(err)
	}

	records := s.Parser.Parse(body)
	if records == nil {
		return nil
	}

	bulk := s.Mongo.RecordsColl.Bulk()
	var wSize int

	for i, record := range *records {
		record.PublisherName = publisherName
		// INDEX adstxtcrawler.records {"publisherName": 1, "domainName": 1, "publisherAccountId": 1, "accountType": 1}
		q := bson.M{"publisherName": publisherName, "domainName": record.DomainName, "publisherAccountId": record.PublisherAccountID, "accountType": record.AccountType}
		u := bson.M{"$set": record}
		bulk.Upsert(q, u)
		wSize++

		if wSize == bulkMaxWriteSizePerBatch || i == len(*records)-1 {
			_, err = bulk.Run()
			if err != nil {
				return err
			}

			bulk = s.Mongo.RecordsColl.Bulk()
			wSize = 0
		}
	}

	return nil
}
