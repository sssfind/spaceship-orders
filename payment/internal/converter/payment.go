package converter

import (
	"payment/internal/model"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

// ToDomainMethod переводит gRPC enum в доменный тип
func ToDomainMethod(pbMethod pb.PaymentMethod) model.PaymentMethod {
	switch pbMethod {
	case pb.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.MethodCard
	case pb.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.MethodSbp
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.MethodCreditCard
	case pb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.MethodInvestorMoney
	default:
		return model.MethodUnspecified
	}
}
