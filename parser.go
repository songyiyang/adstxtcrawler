package adstxtcrawler

import (
	"bytes"
	"log"
	"regexp"
	"strings"
)

const (
	recordSeparatorASCII = byte(10) // \n
	commentSymbolASCII   = byte(35) // #
	fieldSeparatorASCII  = byte(44) // ,
)

var domainRegexp *regexp.Regexp

func init() {
	var err error
	// this regex is referring from https://www.socketloop.com/tutorials/golang-use-regular-expression-to-validate-domain-name
	// currently taking up to third level(e.g. ads.adaptv.advertising.com)
	if domainRegexp, err = regexp.Compile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9])).(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9])).([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}.[a-zA-Z]{2,3})$`); err != nil {
		panic(err)
	}
}

type Parser interface {
	Parse(raw *[]byte) *Records
}

type DefaultParser struct {
	Log *log.Logger
}

var _ Parser = &DefaultParser{}

// Parse will take the raw bytes from http response body and return all the valid iab records
func (s *DefaultParser) Parse(raw *[]byte) *Records {
	if raw == nil {
		return nil
	}

	var records Records
	for _, line := range bytes.Split(*raw, []byte{recordSeparatorASCII}) {
		line = s.trimComments(line)

		rFields := bytes.Split(line, []byte{fieldSeparatorASCII})

		// if record is not valid, try parsing next line
		if !s.isRecordValid(rFields) {
			continue
		}

		record := Record{
			DomainName:         s.getFieldValue(rFields[0]),
			PublisherAccountID: s.getFieldValue(rFields[1]),
			AccountType:        s.getAccountType(rFields[2]),
		}
		if len(rFields) == 4 {
			record.CertificationAuthorityID = s.getFieldValue(rFields[3])
		}

		records = append(records, record)
	}

	return &records
}

func (s *DefaultParser) trimComments(line []byte) []byte {

	if idx := bytes.Index(line, []byte{commentSymbolASCII}); idx != -1 {
		line = line[:idx]
	}

	return line
}

func (s *DefaultParser) getFieldValue(field []byte) string {
	return string(bytes.TrimSpace(field))
}

func (s *DefaultParser) getAccountType(field []byte) AccountType {
	t := s.getFieldValue(field)

	return AccountType(strings.ToUpper(t))
}

// For now, this function only verify the email and account type format.
func (s *DefaultParser) isRecordValid(fields [][]byte) bool {

	if len(fields) != 3 && len(fields) != 4 {
		return false
	}

	domain := s.getFieldValue(fields[0])
	if !domainRegexp.MatchString(domain) {
		return false
	}

	accountType := s.getAccountType(fields[2])
	switch accountType {
	case AccountTypeDirect, AccountTypeReseller:
	default:
		// TODO maybe persist this into another collection (e.g unrendered)?
		s.Log.Printf("Unsupported account type: [%s]", accountType)
		return false
	}

	return true
}
