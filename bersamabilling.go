package bersamabilling

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
    "errors"
    "fmt"
	"io"
	"strings"
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
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = fmt.Sprintf("%s%s", bb.BaseURL, path)
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
	response.XMLName = xml.Name{Local: "return"}

	// set header
	headers := make(map[string]string)
	pathURL := "/api/tfp/generatePaymentCode"

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

	err = bb.call("POST", pathURL, bytes.NewBuffer(payload), &response, headers)
	if err != nil {
		return response, err
	}

	err = bb.GetFinalError(response)
	if err != nil {
	    return response, err
    }

	return response, nil
}

func (bb *BersamaBilling) GetFinalError (payload CreatePaymentCodeResponse) error {
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
