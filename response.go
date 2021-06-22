package bersamabilling

import "encoding/xml"

type CreatePaymentCodeResponse struct {
	XMLName   xml.Name
	Type      string `xml:"type"`
	Ack       string `xml:"ack"`
	BookingID string `xml:"bookingid"`
	VaID      string `xml:"vaid"`
	BankCode  string `xml:"bankcode"`
	Signature string `xml:"signature"`
}

type StatusInquiryResponse struct {
	XMLName xml.Name
	Type    string         `xml:"type"`
	Item    []ItemResponse `xml:",any"`
}

type ItemResponse struct {
	XMLName       xml.Name
	BookingID     string `xml:"bookingid"`
	VaID          string `xml:"vaid"`
	CliendID      string `xml:"clientid"`
	CustomerName  string `xml:"customer_name"`
	IssuerBank    string `xml:"issuer_bank"`
	IssuerName    string `xml:"issuer_name"`
	Amount        string `xml:"amount"`
	ProductID     string `xml:"productid"`
	TrxID         string `xml:"trxid"`
	TrxDate       string `xml:"trx_date"`
	Status        string `xml:"status"`
	StatusMessage string `xml:"-"`
	Signature     string `xml:"signature"`
}
