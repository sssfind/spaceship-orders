package app

import (
	paymentV1 "payment/internal/api/payment/v1"
	"payment/internal/config"
	"payment/internal/service"
	paymentService "payment/internal/service/payment"
	pb "spaceship-orders/shared/pkg/proto/payment/v1"
)

type serviceProvider struct {
	cfg        *config.Config
	paymentSrv service.PaymentService
	apiHandler pb.PaymentServiceServer
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (sp *serviceProvider) PaymentService() (service.PaymentService, error) {
	if sp.paymentSrv == nil {
		sp.paymentSrv = paymentService.NewService()
	}
	return sp.paymentSrv, nil
}

func (sp *serviceProvider) APIHandler() (pb.PaymentServiceServer, error) {
	if sp.apiHandler == nil {
		srv, err := sp.PaymentService()
		if err != nil {
			return nil, err
		}
		sp.apiHandler = paymentV1.NewAPI(srv)
	}
	return sp.apiHandler, nil
}
