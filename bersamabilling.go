package bersamabilling

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"time"
)

type BersamaBilling struct {
	BaseURL  string
	Username string
	Password string
	client   *bersamaBillingHttp
	*logger
}

func New(baseUrl string, username string, password string, timeout time.Duration) *BersamaBilling {
	httpRequest := newRequest(timeout)
	return &BersamaBilling{
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
		client:   httpRequest,
		logger:   newLogger(),
	}
}

func (bb *BersamaBilling) call(method string, path string, body io.Reader, v interface{}, headers map[string]string) error {
	bb.info().Println("Starting http call..")

	return bb.client.exec(method, path, body, v, headers)
}

func (bb *BersamaBilling) CreatePaymentCode(request CreatePaymentCodeRequest) (resp CreatePaymentCodeResponse, err error) {
	bb.info().Println("Starting Get Checkout URL Bersama Billing")
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if err != nil {
			bb.error().Println(err.Error())
		}
	}()
	var response CreatePaymentCodeResponse
	response.XMLName = xml.Name{Local: "Return"}

	// set header
	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	headers["Accept"] = "application/xml"

	// create signature
	signature := fmt.Sprintf("%s%s%s", bb.Username, bb.Password, request.BookingID)
	hash := md5.Sum([]byte(signature))
	request.XMLName = xml.Name{Local: "data"}
	request.Username = bb.Username
	request.Signature = hex.EncodeToString(hash[:])
	//Marshal Order
	payload, errPayload := xml.Marshal(request)
	if errPayload != nil {
		return response, err
	}

	err = bb.call("POST", request.URLCheckout, bytes.NewBuffer(payload), &response, headers)
	if err != nil {
		return response, err
	}

	err = bb.GetFinalError(response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (bb *BersamaBilling) StatusInquiryPayment(request StatusInquiryPaymentRequest) (resp StatusInquiryResponse, err error) {
	bb.info().Println("Starting Get Status Inquiry Bersama Billing")
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if err != nil {
			bb.error().Println(err.Error())
		}
	}()
	var response StatusInquiryResponse
	response.XMLName = xml.Name{Local: "return"}

	// set header
	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	headers["Accept"] = "application/xml"

	//Inject Signature and Username
	for i, v := range request.Item {
		// create signature
		signature := fmt.Sprintf("%s%s%s", bb.Username, bb.Password, v.BookingID)
		hash := md5.Sum([]byte(signature))

		indexData := fmt.Sprintf("%d", i+1)
		if i+1 < 10 {
			indexData = fmt.Sprintf("0%d", i+1)
		}

		xmlName := fmt.Sprintf("item%s", indexData)
		request.Item[i].XMLName = xml.Name{Local: xmlName}
		request.Item[i].Username = bb.Username
		request.Item[i].Signature = hex.EncodeToString(hash[:])
	}

	//Marshal Order
	payload, errPayload := xml.Marshal(request)
	if errPayload != nil {
		return response, err
	}

	err = bb.call("POST", request.URLGetStatus, bytes.NewBuffer(payload), &response, headers)
	if err != nil {
		return response, err
	}

	for _, v := range response.Item {
		response.Item[0].StatusMessage = bb.GetFinalStatusInquiry(v.Status)
	}

	return response, nil
}

func (bb *BersamaBilling) GetFinalError(payload CreatePaymentCodeResponse) error {
	switch payload.Ack {
	case "01":
		return errors.New("Illegal Signature")
	case "02":
		return errors.New("Error Tag XML")
	case "03":
		return errors.New("Error Content Type")
	case "04":
		return errors.New("Error Content Length")
	case "05":
		return errors.New("Error Accessing Database")
	}
	return nil
}

func (bb *BersamaBilling) GetFinalStatusInquiry(code string) string {
	output := ""
	switch code {
	case "00":
		output = "Transaction Success"
		break
	case "05":
		output = "Transaction failed / general error"
		break
	case "13":
		output = "INVALID AMOUNT"
		break
	case "78":
		output = "ALREADY PAID (Transaction Success)"
		break
	case "76":
		output = "Maturity time expired"
		break
	case "101":
		output = "Ilegal signature (inquiry status)"
		break
	case "312":
		output = "Empty Stock (Optional)"
		break
	case "300":
		output = "No Transaction Found"
		break
	case "400":
		output = "Success Transfer, REVERSAL Rejected (Transaction Success)"
		break
	case "401":
		output = "Transfer Cancel, REVERSAL Approve"
		break
	}
	return output
}
