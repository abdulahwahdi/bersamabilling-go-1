package main

import (
	"fmt"
	"github.com/abdulahwahdi/bersamabilling-go"
	"time"
)

const (
	Username = "bene.prabowo@bhinneka.com"
	Password = "PrmL7IvB"
	BaseURL  = "https://bersamabilling.id/portal/index.php"
	Timeout  = 60 * time.Second
)

func main() {
	newBB := bersamabilling.New(BaseURL, Username, Password, Timeout)
	req := bersamabilling.CreatePaymentCodeRequest{
		Type:            "reqpaymentcode",
		BookingID:       "20112345614",
		ClientID:        "1",
		CustomerName:    "Test Data",
		Amount:          10000,
		ProductID:       "01",
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
