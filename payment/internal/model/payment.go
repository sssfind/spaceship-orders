package model

type PaymentMethod string

const (
	MethodUnspecified   PaymentMethod = "UNSPECIFIED"
	MethodCard          PaymentMethod = "CARD"
	MethodSbp           PaymentMethod = "SBP"
	MethodCreditCard    PaymentMethod = "CREDIT_CARD"
	MethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PaymentTransaction struct {
	OrderUUID       string
	UserUUID        string
	TransactionUUID string
	Method          PaymentMethod
}
