package adstxtcrawler

const (
	AccountTypeDirect   AccountType = "DIRECT"
	AccountTypeReseller AccountType = "RESELLER"
)

type AccountType string

type Record struct {
	// PublisherName is the name of the given publisher.
	PublisherName string `bson:"publisherName"`
	// DomainName is the domain name of the advertising system.
	DomainName string `bson:"domainName"`
	// AccountID is the identifier associated with the seller or reseller account within the advertising system in DomainName.
	PublisherAccountID string `bson:"publisherAccountId"`
	// AccountType is an enumeration of the type of account. ['DIRECT', 'RESELLER']
	AccountType AccountType `bson:"accountType"`
	// AuthorityID is an ID that uniquely identifies the advertising system within a certification authority.
	CertificationAuthorityID string `bson:"certificationAuthorityId,omitempty"`
}

type Records []Record
