# Golang ads.txt Web Crawler

Crawler for ads.txt files given a list of URLs or domains etc and saves them to a MongoDB.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Dependencies:

```
go
mongodb[3.4.9+]
```

### Installing

OSX:

```
brew update
brew install go
brew install mongodb
```

*NIX and WINDOWS, please follow the instruction links

```
https://golang.org/doc/install
https://docs.mongodb.com/manual/administration/install-on-linux/

```

For sample test: setup GOROOT and GOPATH, then make sure mognodb is running on localhost:27017

## Running the example

Run the example main binary and tests for this system

### Build example binary


```
cd cmd/example
go build *.go
```

### Run example binary

```
./main
```

### Run tests
```
go test
```

### Curl request from HTTP endpoint

For request per publisher
```
curl http://localhost:8000/adstxt?publisherName=Bloomberg
curl http://localhost:8000/adstxt?publisherName=WordPress
```

For request of multi-publishers
```
curl http://localhost:8000/adstxt?publisherName=CNN&publisherName=Gizmodo&publisherName=NYTimes
```

For request with multiple filters
```
curl http://localhost:8000/adstxt?publisherName=CNN&accountType=DIRECT
curl http://localhost:8000/adstxt?publisherName=WordPress&domainName=pubmine%2Ecom&publisherAccountId=3
```

## TODOs
* Handles the case when AccountType is not DIRECT or RESELLER
* Moves tests to use random DB ports instead of using defaults
* Builds higher level system to reduce redundancies of fetching same publisher
* etc

## Authors

* **Yiyang Song** - *Initial work* - [songyiyang](https://github.com/songyiyang)

## Acknowledgments

* Domain matching regex is referring from: https://www.socketloop.com/tutorials/golang-use-regular-expression-to-validate-domain-name
* For more details about ads.txt and IAB, visit https://iabtechlab.com/~iabtec5/wp-content/uploads/2016/07/IABOpenRTBAds.txtSpecification_Version1_Final.pdf
* etc
