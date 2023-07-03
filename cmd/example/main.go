package main

import (
	"fmt"
	"github.com/ysong/adstxtcrawler"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	httpEndpointHost = "0.0.0.0"
	httpEndpointPort = "8000"
)

var logger *log.Logger
var mongo *adstxtcrawler.MongoCxt

func init() {
	logger = log.New(os.Stdout, "[adstxtcrawler]: ", 0)

	var err error
	mongo, err = adstxtcrawler.InitDefaultMongoCxt() // default mongodb will be running on 127.0.0.1:27017
	if err != nil {
		panic(err)
	}

	if err = mongo.DropDB(); err != nil {
		panic(err)
	}
}

func main() {

	wg := sync.WaitGroup{}
	defer wg.Wait()

	crawler := adstxtcrawler.Crawler{&wg, logger, mongo, &adstxtcrawler.DefaultParser{logger}}

	wg.Add(5)
	go crawler.Do("CNN", "http://www.cnn.com/ads.txt", nil)
	go crawler.Do("Gizmodo", "http://www.gizmodo.com/ads.txt", nil)
	go crawler.Do("NYTimes", "http://www.nytimes.com/ads.txt", nil)
	go crawler.Do("Bloomberg", "https://www.bloomberg.com/ads.txt", nil)
	go crawler.Do("WordPress", "https://wordpress.com/ads.txt", nil)

	http.HandleFunc("/adstxt", func(w http.ResponseWriter, r *http.Request) {
		results, err := crawler.RenderQueryResults(r.URL.Query())
		if err != nil {
			http.Error(w, "Unable to process the request at this moment.", http.StatusInternalServerError)
			crawler.Log.Print(err.Error())
			return
		}

		w.Write(*results)
	})

	wg.Add(1)
	go func() {
		logger.Printf("Listening on http://" + httpEndpointHost + ":" + httpEndpointPort + "/")
		defer wg.Done()

		if err := http.ListenAndServe(httpEndpointHost+":"+httpEndpointPort, nil); err != nil {
			fmt.Printf("could not start data endpoint HTTP server")
			return
		}

	}()
}
