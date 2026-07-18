package converter

import (
	"iam/internal/model"
	commonv1 "spaceship-orders/shared/pkg/proto/common/v1"
)

func ToProtoNotificationMethods(methods []model.NotificationMethod) []*commonv1.NotificationMethod {
	if methods == nil {
		return nil
	}
	protoMethods := make([]*commonv1.NotificationMethod, 0, len(methods))
	for _, m := range methods {
		protoMethods = append(protoMethods, &commonv1.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}
	return protoMethods
}

func ToModelNotificationMethods(protoMethods []*commonv1.NotificationMethod) []model.NotificationMethod {
	if protoMethods == nil {
		return nil
	}
	methods := make([]model.NotificationMethod, 0, len(protoMethods))
	for _, m := range protoMethods {
		methods = append(methods, model.NotificationMethod{
			ProviderName: m.GetProviderName(),
			Target:       m.GetTarget(),
		})
	}
	return methods
}

func ToProtoUser(u *model.User) *commonv1.User {
	if u == nil {
		return nil
	}
	return &commonv1.User{
		UserUuid:            u.UUID,
		Login:               u.Login,
		Email:               u.Email,
		NotificationMethods: ToProtoNotificationMethods(u.NotificationMethods),
	}
}
