package converter

import (
	"encoding/json"

	"iam/internal/model"
	repoModel "iam/internal/repository/model"
)

func ToUserFromRepo(rm *repoModel.UserRepoModel) (*model.User, error) {
	var methods []model.NotificationMethod
	if len(rm.NotificationMethods) > 0 {
		if err := json.Unmarshal(rm.NotificationMethods, &methods); err != nil {
			return nil, err
		}
	}

	return &model.User{
		UUID:                rm.UUID,
		Login:               rm.Login,
		PasswordHash:        rm.PasswordHash,
		Email:               rm.Email,
		NotificationMethods: methods,
	}, nil
}
