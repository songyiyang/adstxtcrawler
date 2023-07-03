package adstxtcrawler

import (
	"log"
	"os"
	"testing"
)

func TestDefaultParse(t *testing.T) {

	dp := DefaultParser{log.New(os.Stdout, "", 0)}

	for i, s := range []struct {
		raw      string
		expected Records
	}{
		{"", nil},
		{" # CNN.com/ads.txt \n  # \n # DOMESTIC \n ", nil},
		{"google.com, pub-7439281311086140, DIRECT, f08c47fec0942fa0 # banner, video, native \n", Records{Record{DomainName: "google.com", PublisherAccountID: "pub-7439281311086140", AccountType: AccountTypeDirect, CertificationAuthorityID: "f08c47fec0942fa0"}}},
		{"c.amazon-adsystem.com, 3159, DIRECT # banner, video", Records{Record{DomainName: "c.amazon-adsystem.com", PublisherAccountID: "3159", AccountType: AccountTypeDirect}}},
		{"google.com, # 1425 , DIRECT", nil},
		{"google.com, 1425 # random, DIRECT", nil},
	} {
		rawBytes := []byte(s.raw)
		res := dp.Parse(&rawBytes)

		if len(*res) != len(s.expected) {
			t.Errorf("t%d: want %d records, got %d from Parse call of [%s]", i, len(s.expected), len(*res), s.raw)
		}

		if len(*res) > 0 {
			if (*res)[0].DomainName != s.expected[0].DomainName {
				t.Errorf("t%d: want [%s] as DomainName, got [%s] from Parse call of [%s]", i, s.expected[0].DomainName, (*res)[0].DomainName, s.raw)
			}

			if (*res)[0].AccountType != s.expected[0].AccountType {
				t.Errorf("t%d: want [%s] as AccountType, got [%s] from Parse call of [%s]", i, s.expected[0].AccountType, (*res)[0].AccountType, s.raw)
			}

			if (*res)[0].PublisherAccountID != s.expected[0].PublisherAccountID {
				t.Errorf("t%d: want [%s] as PublisherAccountID, got [%s] from Parse call of [%s]", i, s.expected[0].PublisherAccountID, (*res)[0].PublisherAccountID, s.raw)
			}

			if (*res)[0].CertificationAuthorityID != s.expected[0].CertificationAuthorityID {
				t.Errorf("t%d: want [%s] as CertificationAuthorityID, got [%s] from Parse call of [%s]", i, s.expected[0].CertificationAuthorityID, (*res)[0].CertificationAuthorityID, s.raw)
			}
		}
	}
}
