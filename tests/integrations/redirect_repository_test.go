package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/domain"
	redirectRepo "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/repository/postgresql"
)

func TestRedirectRepository(t *testing.T) {
	if testPool == nil {
		t.Skip("Skipping integration test: testPool is nil")
	}
	ctx := context.Background()
	repo := redirectRepo.NewRepository(testPool)

	var instID string
	err := testPool.QueryRow(ctx, "SELECT id FROM auth.institutions WHERE code = 'TECH-UNI'").Scan(&instID)
	require.NoError(t, err)

	t.Run("FindUniversityById", func(t *testing.T) {
		inst, err := repo.FindUniversityById(ctx, instID)
		assert.NoError(t, err)
		assert.Equal(t, instID, inst.Id)
	})

	t.Run("FindUserBySub", func(t *testing.T) {
		sub := "sub_admin_tech_001"
		user, err := repo.FindUserBySub(ctx, instID, sub)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, sub, user.ExternalSubject)
		// Admin seeded with roles, check if loaded
		assert.NotEmpty(t, user.Roles)
	})

	t.Run("StoreSession", func(t *testing.T) {
		// Needs seed user
		var userID string
		err := testPool.QueryRow(ctx, "SELECT id FROM auth.users WHERE external_subject = 'sub_admin_tech_001'").Scan(&userID)
		require.NoError(t, err)

		session := &domain.Session{
			SessionId:       "sess_test_001",
			UniversityId:    instID,
			UserId:          userID,
			ExternalSubject: "sub_admin_tech_001",
			Roles:           []string{"admin"},
			AccessToken:     "access_token_123",
			ExpiresAt:       time.Now().Add(1 * time.Hour),
		}

		err = repo.StoreSession(ctx, session)
		assert.NoError(t, err)

		// Verification verify via SQL as repo doesn't have FindSession?
		// Actually Middleware uses auth.sessions, so it's critical.
		var dbSession string
		err = testPool.QueryRow(ctx, "SELECT session_id FROM auth.sessions WHERE session_id = 'sess_test_001'").Scan(&dbSession)
		assert.NoError(t, err)
		assert.Equal(t, "sess_test_001", dbSession)
	})
	t.Run("FindInstitutionByID_NotFound", func(t *testing.T) {
		inst, err := repo.FindUniversityById(ctx, "00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Nil(t, inst)
	})

	t.Run("FindUserBySub_NotFound", func(t *testing.T) {
		user, err := repo.FindUserBySub(ctx, instID, "non-existent")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	t.Run("StoreSession_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := repo.StoreSession(ctx, &domain.Session{})
		assert.Error(t, err)
		t.Run("FindInstitutionByID_ContextCancel", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := repo.FindUniversityById(ctx, "some-id")
			assert.Error(t, err)
		})

		t.Run("FindUserBySub_ContextCancel", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := repo.FindUserBySub(ctx, instID, "some-sub")
			assert.Error(t, err)
		})
	})
}
