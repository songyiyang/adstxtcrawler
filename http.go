package adstxtcrawler

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (s *Crawler) getURLBody(url string, headers map[string]string) (*[]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return s.parseValidHTTPResponseBody(response)
}

func (s *Crawler) parseValidHTTPResponseBody(response *http.Response) (*[]byte, error) {
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if body, err := ioutil.ReadAll(response.Body); err != nil {
			return nil, err
		} else {
			return &body, nil
		}

	} else {
		return nil, errors.New(fmt.Sprintf("ERROR status code for %s: %s", response.Request.URL, response.Status))
	}
}

// RenderQueryResults will parse http endpoint query and render results from db
func (s *Crawler) RenderQueryResults(query url.Values) (ret *[]byte, err error) {
	q := make(bson.M)

	for k, v := range query {
		q[k] = bson.M{"$in": v}
	}

	iter := s.Mongo.RecordsColl.Find(q).Iter()
	defer func() {
		err = iter.Close()
	}()

	var record Record
	var line string
	ret = &[]byte{}
	for iter.Next(&record) {
		line = fmt.Sprintf("%s,%s,%v", record.DomainName, record.PublisherAccountID, record.AccountType)
		if record.CertificationAuthorityID != "" {
			line += fmt.Sprintf(",%s", record.CertificationAuthorityID)
		}
		line += "\n"

		*ret = append(*ret, []byte(line)...)
	}

	return
}
