package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

type MockSessionConfig struct {
	mock.Mock
}

func (m *MockSessionConfig) TTL() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (s *AuthServiceTestSuite) TestLogin_Success() {
	login := "darth_vader"
	password := "deathstar123"
	ttl := 24 * time.Hour

	// Генерируем валидный bcrypt-хэш для пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	s.NoError(err)

	dbUser := &model.User{
		UUID:         uuid.New().String(),
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	cfgMock := new(MockSessionConfig)
	s.service = NewAuthService(s.userRepo, s.sessionRepo, cfgMock)

	s.userRepo.On("GetByLogin", s.ctx, login).Return(dbUser, nil).Once()
	cfgMock.On("TTL").Return(ttl).Once()

	// Используем MatchedBy, так как внутри генерируется случайный UUID сессии и текущее время
	s.sessionRepo.On("Create", s.ctx, mock.MatchedBy(func(session *model.Session) bool {
		return session.UserUUID == dbUser.UUID && session.SessionUUID != ""
	}), ttl).Return(nil).Once()

	sessionUUID, err := s.service.Login(s.ctx, login, password)

	s.NoError(err)
	s.NotEmpty(sessionUUID)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
	cfgMock.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_UserNotFound() {
	login := "non_existent_user"
	password := "any_password"

	// Эмулируем ошибку из базы данных (пользователь не найден)
	s.userRepo.On("GetByLogin", s.ctx, login).
		Return(nil, errors.New("user not found in db")).Once()

	sessionUUID, err := s.service.Login(s.ctx, login, password)

	// Проверяем, что метод скрывает системную ошибку и возвращает ErrInvalidCredentials
	s.Empty(sessionUUID)
	s.ErrorIs(err, model.ErrInvalidCredentials)
	s.userRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_WrongPassword() {
	login := "luke_skywalker"
	correctPassword := "i_am_your_father"
	wrongPassword := "wrong_pass"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	s.NoError(err)

	dbUser := &model.User{
		UUID:         uuid.New().String(),
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	s.userRepo.On("GetByLogin", s.ctx, login).Return(dbUser, nil).Once()

	// Передаем неверный пароль
	sessionUUID, err := s.service.Login(s.ctx, login, wrongPassword)

	s.Empty(sessionUUID)
	s.ErrorIs(err, model.ErrInvalidCredentials)
	s.userRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_SessionCreateError() {
	login := "obi_wan"
	password := "high_ground"
	ttl := time.Hour

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	s.NoError(err)

	dbUser := &model.User{
		UUID:         uuid.New().String(),
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	cfgMock := new(MockSessionConfig)
	s.service = NewAuthService(s.userRepo, s.sessionRepo, cfgMock)

	dbErr := errors.New("redis cluster is down")

	s.userRepo.On("GetByLogin", s.ctx, login).Return(dbUser, nil).Once()
	cfgMock.On("TTL").Return(ttl).Once()
	s.sessionRepo.On("Create", s.ctx, mock.Anything, ttl).Return(dbErr).Once()

	sessionUUID, err := s.service.Login(s.ctx, login, password)

	s.Empty(sessionUUID)
	s.Error(err)
	s.ErrorIs(err, dbErr)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}
