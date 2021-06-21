package bersamabilling

import "encoding/xml"

type CreatePaymentCodeRequest struct {
	XMLName         xml.Name
	Type            string `xml:"type"`
	BookingID       string `xml:"bookingid"`
	ClientID        string `xml:"clientid"`
	CustomerName    string `xml:"customer_name"`
	Amount          int64  `xml:"amount"`
	ProductID       string `xml:"productid"`
	Interval        int    `xml:"interval"`
	Username        string `xml:"username"`
	BookingDatetime string `xml:"booking_datetime"`
	Signature       string `xml:"signature"`
}
