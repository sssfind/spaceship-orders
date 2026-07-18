package converter

import (
	"iam/internal/model"
	repoModel "iam/internal/repository/model"
	"time"
)

func ToRepoFromSession(m *model.Session) *repoModel.SessionRepoModel {
	return &repoModel.SessionRepoModel{
		SessionUUID: m.SessionUUID,
		UserUUID:    m.UserUUID,
		CreatedAt:   m.CreatedAt.Unix(),
	}
}

func ToSessionFromRepo(rm *repoModel.SessionRepoModel) *model.Session {
	return &model.Session{
		SessionUUID: rm.SessionUUID,
		UserUUID:    rm.UserUUID,
		CreatedAt:   time.Unix(rm.CreatedAt, 0),
	}
}
