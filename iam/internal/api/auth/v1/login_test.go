// iam/internal/api/auth/v1/login_test.go
package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	authv1 "spaceship-orders/shared/pkg/proto/auth/v1"
)

func (s *LoginTestSuite) TestLogin_Success() {
	ctx := context.Background()
	login := "space_cadet"
	password := "secure_pass"
	expectedSessionUUID := uuid.New().String()

	req := &authv1.LoginRequest{
		Login:    login,
		Password: password,
	}

	s.authServiceMock.On("Login", ctx, login, password).
		Return(expectedSessionUUID, nil).Once()

	resp, err := s.api.Login(ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(expectedSessionUUID, resp.SessionUuid)
	s.authServiceMock.AssertExpectations(s.T())
}

func (s *LoginTestSuite) TestLogin_InvalidCredentials() {
	ctx := context.Background()
	login := "imposter"
	password := "wrong_pass"

	req := &authv1.LoginRequest{
		Login:    login,
		Password: password,
	}

	// Эмулируем ошибку валидации учетных данных из бизнес-логики
	s.authServiceMock.On("Login", ctx, login, password).
		Return("", model.ErrInvalidCredentials).Once()

	resp, err := s.api.Login(ctx, req)

	// Проверяем gRPC status-код ошибки
	s.Nil(resp)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Unauthenticated, st.Code()) // Хендлер должен вернуть 16 (Unauthenticated)
	s.Equal(model.ErrInvalidCredentials.Error(), st.Message())
	s.authServiceMock.AssertExpectations(s.T())
}

func (s *LoginTestSuite) TestLogin_InternalError() {
	ctx := context.Background()
	req := &authv1.LoginRequest{
		Login:    "error_user",
		Password: "any_password",
	}

	dbErr := errors.New("database connection lost")

	// Эмулируем непредвиденную инфраструктурную ошибку
	s.authServiceMock.On("Login", ctx, mock.Anything, mock.Anything).
		Return("", dbErr).Once()

	resp, err := s.api.Login(ctx, req)

	// Проверяем маппинг в Internal ошибку
	s.Nil(resp)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Internal, st.Code()) // Хендлер должен скрыть детали ошибки и вернуть Internal
	s.Equal("failed to login", st.Message())
	s.authServiceMock.AssertExpectations(s.T())
}
