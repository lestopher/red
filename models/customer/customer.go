package customer

import (
	"errors"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErpCodes contains a title and where the moveTo should go
type ErpCodes struct {
	title  string
	moveTo string
}

var (
	// UnverifiedEducationErpCodes a map of unverified erp codes and to what they are mapped
	UnverifiedEducationErpCodes = map[string]ErpCodes{
		"ESV": ErpCodes{title: "Student", moveTo: "ESN"},
		"EIV": ErpCodes{title: "Institution", moveTo: "EIV"},
		"ETN": ErpCodes{title: "Faculty Member", moveTo: "ETV"},
	}
	// VerifiedEducationErpCodes a map of verified erp codes and to what they are mapped
	VerifiedEducationErpCodes = map[string]ErpCodes{
		"ESN": ErpCodes{title: "Student", moveTo: "ESV"},
		"EIV": ErpCodes{title: "Institution", moveTo: "EIN"},
		"ETV": ErpCodes{title: "Faculty Member", moveTo: "ETN"},
	}
	discountGroupCodes = map[string]struct{}{
		"ESV": struct{}{},
		"ETV": struct{}{},
		"EIV": struct{}{},
	}
	// DBSession holds a session
	DBSession *Session
)

// Session represents a mongo session
type Session struct {
	*mgo.Session
}

// NewSession provides a pointer to a Session and an error response
func NewSession() (*Session, error) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	return &Session{session}, err
}

// FindByID returns a pointer to a customer record
func FindByID(id string) (*Customer, error) {
	if DBSession == nil {
		return nil, errors.New("Database session is nil")
	}

	var c Customer
	err := DBSession.DB("red_development").C("customers").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Customer is mapped to an object in mongodb
type Customer struct {
	// Fields would be equivalent to attr_accessor in rails
	documentHash           string
	organizationRecID      string
	organizationRecVersion string
	dirPartyRecID          string
	dirPartyRecVersion     string
	custGroup              string
	paymTermID             string

	// Fields would be equivalent to attr_accessible in rails
	ID                  bson.ObjectId `bson:"_id,omitempty"`
	Name                string        `bson:"name"`
	AccountNumber       string        `bson:"account_number"`
	AccountBalance      float32       `bson:"account_balance"`
	ErpCode             string        `bson:"erp_code"`
	CountryCode         string        `bson:"country_code"`
	DefaultShippingID   bson.ObjectId `bson:"default_shipping_id"`
	DefaultBillingID    bson.ObjectId `bson:"default_billing_id"`
	DownloadPermissions []interface{} `bson:"download_permissions"`
	Blocked             bool          `bson:"blocked"`
	RedOneShipped       bool          `bson:"red_one_shipped"`
	LegacyID            int32         `bson:"legacy_id"`
	PinNumbers          []int32       `bson:"pin_numbers"`
	EarlyAdopter        bool          `bson:"early_adopter"`
	EducationExpiresAt  time.Time     `bson:"education_expires_at"`
	VatNumber           string        `bson:"vat_numbers"`
	ValidatedVatName    string        `bson:"validated_vat_name"`
	Multiple            bool          `bson:"multiple"`
}

// NewCustomer creates a new Customer
func NewCustomer() *Customer {
	return &Customer{
		DownloadPermissions: make([]interface{}, 0),
		Blocked:             false,
		RedOneShipped:       false,
		PinNumbers:          make([]int32, 0),
		EarlyAdopter:        false,
		Multiple:            false,
	}
}

// IsEducationDiscountApproved returns a bool as to whether the discount is approved
func (c Customer) IsEducationDiscountApproved() bool {
	_, ok := VerifiedEducationErpCodes[c.ErpCode]
	return ok
}

// IsEducationDiscountDenied returns a bool as to whether the discount is denied
func (c Customer) IsEducationDiscountDenied() bool {
	_, ok := UnverifiedEducationErpCodes[c.ErpCode]
	return ok
}

// IsEducationDiscountVerified returns a bool as to whether the discount is verified
func (c Customer) IsEducationDiscountVerified() bool {
	_, ok := VerifiedEducationErpCodes[c.ErpCode]
	return ok
}

// IsEducationDiscountUnverified returns a bool as to whether the discount is unverified
func (c Customer) IsEducationDiscountUnverified() bool {
	_, ok := UnverifiedEducationErpCodes[c.ErpCode]
	return ok
}

// HasEducationProgram returns as bool as to whether the customer has an education program
func (c Customer) HasEducationProgram() bool {
	return c.IsEducationDiscountUnverified() || c.IsEducationDiscountVerified()
}

// EducationDiscountStatus gives the status or an error depending on if the customer
// has an education program or not
func (c Customer) EducationDiscountStatus() (string, error) {
	if !c.HasEducationProgram() {
		return "", errors.New("Customer does not have an education discount.")
	}
	switch {
	case c.IsEducationDiscountVerified():
		return "pending", nil
	case c.IsEducationDiscountUnverified():
		return "not verified", nil
	case c.IsEducationDiscountVerified():
		return "verified", nil
	case c.IsEducationDiscountUnverified():
		return "not verified", nil
	default:
		return "", errors.New("Didn't match any case statements.")
	}
}

func (c Customer) isEducationDiscountExpired() bool {
	return c.EducationExpiresAt.Before(time.Now())
}

func (c Customer) isInDiscountGroup() bool {
	if _, ok := discountGroupCodes[c.ErpCode]; ok {
		return c.EducationExpiresAt.After(time.Now())
	}

	return false
}
