package bersamabilling

type BersamaBillingService interface {
	CreatePaymentCode(request CreatePaymentCodeRequest) (resp CreatePaymentCodeResponse, err error)
}
