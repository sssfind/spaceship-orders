package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"iam/internal/model"
	commonv1 "spaceship-orders/shared/pkg/proto/common/v1"
)

func TestToProtoUser_Success(t *testing.T) {
	u := &model.User{
		UUID:  "user-uuid-12345",
		Login: "space_captain",
		Email: "captain@galaxy.com",
		NotificationMethods: []model.NotificationMethod{
			{
				ProviderName: "telegram",
				Target:       "@captain_bot",
			},
		},
	}

	protoUser := ToProtoUser(u)

	assert.NotNil(t, protoUser)
	assert.Equal(t, u.UUID, protoUser.UserUuid)
	assert.Equal(t, u.Login, protoUser.Login)
	assert.Equal(t, u.Email, protoUser.Email)

	assert.Len(t, protoUser.NotificationMethods, 1)
	assert.Equal(t, "telegram", protoUser.NotificationMethods[0].ProviderName)
	assert.Equal(t, "@captain_bot", protoUser.NotificationMethods[0].Target)
}

func TestToProtoUser_Nil(t *testing.T) {
	protoUser := ToProtoUser(nil)
	assert.Nil(t, protoUser)
}

func TestToModelNotificationMethods_Success(t *testing.T) {
	protoMethods := []*commonv1.NotificationMethod{
		{
			ProviderName: "email",
			Target:       "test@test.com",
		},
	}

	models := ToModelNotificationMethods(protoMethods)

	assert.Len(t, models, 1)
	assert.Equal(t, "email", models[0].ProviderName)
	assert.Equal(t, "test@test.com", models[0].Target)
}

func TestToModelNotificationMethods_Nil(t *testing.T) {
	assert.Nil(t, ToModelNotificationMethods(nil))
	assert.Nil(t, ToProtoNotificationMethods(nil))
}
