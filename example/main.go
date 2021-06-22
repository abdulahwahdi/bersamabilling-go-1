package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Bhinneka/bersamabilling-go"
	"time"
)

const (
	Username = ""
	Password = ""
	BaseURL  = "https://bersamabilling.id/portal/index.php"
	Timeout  = 60 * time.Second
)

func main() {
	newBB := bersamabilling.New(BaseURL, Username, Password, Timeout)
	GetPaymentCode(newBB)
	StatusInquiry(newBB)
}

func GetPaymentCode(newBB bersamabilling.BersamaBillingService) {
	req := bersamabilling.CreatePaymentCodeRequest{
		Type:            "reqpaymentcode",
		BookingID:       "20112345614",
		ClientID:        "1",
		CustomerName:    "Test Data",
		Amount:          10000,
		ProductID:       "PDB01",
		Interval:        100,
		BookingDatetime: "2021-06-16 20:20:33",
	}
	resp, err := newBB.CreatePaymentCode(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}

func StatusInquiry(newBB bersamabilling.BersamaBillingService) {
	req := bersamabilling.StatusInquiryPaymentRequest{
		XMLName: xml.Name{Local: "notification"},
		Type:    "reqtrxstatus",
		Item: []bersamabilling.ItemRequest{
			{
				XMLName:         xml.Name{Local: "item01"},
				BookingID:       "20112345614",
				VaID:            "20112345614",
				BookingDatetime: "2021-06-16 20:20:33",
			},
		},
	}
	resp, err := newBB.StatusInquiryPayment(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}
