package v1

import (
	"context"
	"errors"

	"iam/internal/converter"
	"iam/internal/model"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *LoginTestSuite) TestWhoami_Success() {
	ctx := context.Background()
	sessionUUID := uuid.New().String()
	userUUID := uuid.New().String()

	existingUser := &model.User{
		UUID:  userUUID,
		Login: "star_lord",
		Email: "peter@guardians.galaxy",
	}

	req := &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	}

	s.authServiceMock.On("Whoami", ctx, sessionUUID).
		Return(existingUser, nil).Once()

	resp, err := s.api.Whoami(ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(converter.ToProtoUser(existingUser), resp.User)
	s.authServiceMock.AssertExpectations(s.T())
}

func (s *LoginTestSuite) TestWhoami_SessionNotFound() {
	ctx := context.Background()
	sessionUUID := uuid.New().String()

	req := &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	}

	// Эмулируем ошибку отсутствия сессии
	s.authServiceMock.On("Whoami", ctx, sessionUUID).
		Return(nil, model.ErrSessionNotFound).Once()

	resp, err := s.api.Whoami(ctx, req)

	s.Nil(resp)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Unauthenticated, st.Code())
	s.Equal("session expired or invalid", st.Message())
	s.authServiceMock.AssertExpectations(s.T())
}

func (s *LoginTestSuite) TestWhoami_UserNotFound() {
	ctx := context.Background()
	sessionUUID := uuid.New().String()

	req := &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	}

	// Эмулируем ошибку, когда сессия есть, но юзера в базе нет
	s.authServiceMock.On("Whoami", ctx, sessionUUID).
		Return(nil, model.ErrUserNotFound).Once()

	resp, err := s.api.Whoami(ctx, req)

	s.Nil(resp)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.NotFound, st.Code())
	s.Equal("user not found", st.Message())
	s.authServiceMock.AssertExpectations(s.T())
}

func (s *LoginTestSuite) TestWhoami_InternalError() {
	ctx := context.Background()
	sessionUUID := uuid.New().String()

	req := &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	}

	s.authServiceMock.On("Whoami", ctx, mock.Anything).
		Return(nil, errors.New("random internal error")).Once()

	resp, err := s.api.Whoami(ctx, req)

	s.Nil(resp)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Internal, st.Code())
	s.Equal("failed to recognize session", st.Message())
	s.authServiceMock.AssertExpectations(s.T())
}
