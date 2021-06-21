package bersamabilling

import "encoding/xml"

type CreatePaymentCodeResponse struct {
	XMLName xml.Name
	Type      string `xml:"type"`
	Ack       string `xml:"ack"`
	BookingID string `xml:"bookingid"`
	VaID      string `xml:"vaid"`
	BankCode  string `xml:"bankcode"`
	Signature string `xml:"signature"`
}
