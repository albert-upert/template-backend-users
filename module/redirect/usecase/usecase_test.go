package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/albert-upert/template-backend-users/config"
	"github.com/albert-upert/template-backend-users/module/redirect/domain"
	"github.com/albert-upert/template-backend-users/module/redirect/usecase"
	"github.com/albert-upert/template-backend-users/tests/mocks"
	"github.com/albert-upert/template-backend-utils-libraries/idp"
	"github.com/albert-upert/template-backend-utils-libraries/idp/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Redirect(t *testing.T) {
	mockRepo := new(mocks.RedirectRepositoryMock)
	mockIDPProvider := new(mocks.IDPProviderMock)
	mockIDPClient := new(mocks.IDPClientMock)

	conf := &config.InternalAppConfig{
		RedirectUrl: "http://default.com",
	}

	uc := usecase.NewUseCase(conf, mockRepo, mockIDPProvider)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-1"
		token := "valid-token"

		inst := &domain.University{
			Id: instId,
			Settings: idp.Setting{
				Url: "http://inst.com",
			},
		}

		authSession := &client.AuthSession{
			Sub:       "user-sub",
			ExpiresIn: 3600,
		}

		user := &domain.User{
			Id:              "user-1",
			ExternalSubject: "user-sub",
			Roles:           []string{"admin"},
		}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return(mockIDPClient, nil).Once()
		mockIDPClient.On("Check", mock.Anything, token).Return(authSession, nil).Once()
		mockRepo.On("FindUserBySub", mock.Anything, instId, authSession.Sub).Return(user, nil).Once()
		mockRepo.On("StoreSession", mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
			return s.UserId == user.Id && s.AccessToken == token
		})).Return(nil).Once()

		url, sid, err := uc.Redirect(ctx, instId, token)

		assert.NoError(t, err)
		assert.Equal(t, "http://inst.com", url)
		assert.NotEmpty(t, sid)

		mockRepo.AssertExpectations(t)
		mockIDPProvider.AssertExpectations(t)
		mockIDPClient.AssertExpectations(t)
	})

	t.Run("InstitutionNotFound", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-2"
		token := "token"

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return((*domain.University)(nil), errors.New("not found")).Once()

		url, sid, err := uc.Redirect(ctx, instId, token)

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Empty(t, sid)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetIDPFail", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-3"
		token := "token"

		inst := &domain.University{Id: instId}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return((*mocks.IDPClientMock)(nil), errors.New("fail")).Once()

		url, sid, err := uc.Redirect(ctx, instId, token)

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Empty(t, sid)

		mockRepo.AssertExpectations(t)
		mockIDPProvider.AssertExpectations(t)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-4"
		token := "invalid"

		inst := &domain.University{Id: instId}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return(mockIDPClient, nil).Once()
		mockIDPClient.On("Check", mock.Anything, token).Return((*client.AuthSession)(nil), errors.New("invalid")).Once()

		url, _, err := uc.Redirect(ctx, instId, token)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		assert.Empty(t, url)

		mockRepo.AssertExpectations(t)
		mockIDPProvider.AssertExpectations(t)
		mockIDPClient.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-5"
		token := "valid"

		inst := &domain.University{Id: instId}
		authSession := &client.AuthSession{Sub: "sub", ExpiresIn: 60}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return(mockIDPClient, nil).Once()
		mockIDPClient.On("Check", mock.Anything, token).Return(authSession, nil).Once()
		mockRepo.On("FindUserBySub", mock.Anything, instId, authSession.Sub).Return((*domain.User)(nil), errors.New("user missing")).Once()

		url, sid, err := uc.Redirect(ctx, instId, token)

		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Empty(t, sid)

		mockRepo.AssertExpectations(t)
	})

	t.Run("StoreSessionFail", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-6"
		token := "valid"

		inst := &domain.University{Id: instId}
		authSession := &client.AuthSession{Sub: "sub", ExpiresIn: 60}
		user := &domain.User{Id: "u1"}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return(mockIDPClient, nil).Once()
		mockIDPClient.On("Check", mock.Anything, token).Return(authSession, nil).Once()
		mockRepo.On("FindUserBySub", mock.Anything, instId, authSession.Sub).Return(user, nil).Once()
		mockRepo.On("StoreSession", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()

		url, _, err := uc.Redirect(ctx, instId, token)

		assert.Error(t, err)
		assert.Empty(t, url)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DefaultURL", func(t *testing.T) {
		ctx := context.Background()
		instId := "inst-7"
		token := "valid"

		inst := &domain.University{Id: instId, Settings: idp.Setting{Url: ""}} // Empty URL
		authSession := &client.AuthSession{Sub: "sub", ExpiresIn: 60}
		user := &domain.User{Id: "u1"}

		mockRepo.On("FindUniversityById", mock.Anything, instId).Return(inst, nil).Once()
		mockIDPProvider.On("GetIDP", mock.Anything, instId).Return(mockIDPClient, nil).Once()
		mockIDPClient.On("Check", mock.Anything, token).Return(authSession, nil).Once()
		mockRepo.On("FindUserBySub", mock.Anything, instId, authSession.Sub).Return(user, nil).Once()
		mockRepo.On("StoreSession", mock.Anything, mock.Anything).Return(nil).Once()

		url, _, err := uc.Redirect(ctx, instId, token)

		assert.NoError(t, err)
		assert.Equal(t, "http://default.com", url)

		mockRepo.AssertExpectations(t)
	})
}
