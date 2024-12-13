package model
type EmailOrder struct {
	ID       []byte
	HexID    string
	Customer string
	Products []struct {
		ID               []byte
		HexIdProduct     string
		Quantity         uint
		Typ SubscriptionType
		ProductName      string
		ImgUrl           string
	}
	CreateAt     uint
	CreateAtDate string
	ShipInfo     struct {
		FirstName       string
		LastName        string
		Email           string
		Country         string
		City            string
		StateOrProvince string
		PostalCode      string
		Phone           string
		AddressDetail   string
	}
	ShippingFee uint
}
type Data struct {
	Order        EmailOrder
	PaymentOrder uint
}
type Sender struct {
	Address string
	Subject string
}
type Recipient struct {
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
}
type SubscriptionType uint

const (
	None SubscriptionType = iota
    ThreeMonth 
    SixMonth
    TwelveMonth
)