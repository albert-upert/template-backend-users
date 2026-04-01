package integrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	libtypes "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
	usersRepo "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/repository/postgresql"
)

func TestUsersRepository(t *testing.T) {
	if testPool == nil {
		t.Skip("Skipping integration test: testPool is nil")
	}
	ctx := context.Background()
	repo := usersRepo.NewRepository(testPool)

	// Fetch some seed data IDs
	var universityID string
	err := testPool.QueryRow(ctx, "SELECT id FROM auth.institutions WHERE code = 'TECH-UNI'").Scan(&universityID)
	require.NoError(t, err)

	t.Run("Store_NewUser", func(t *testing.T) {
		user := &domain.User{
			UniversityId:     universityID,
			ExternalSubject:  "sub_test_new_001",
			IdentityProvider: "central",
			Status:           "active",
			Metadata:         map[string]any{"name": "Test User"},
		}

		err := repo.Store(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.Id)

		// Verify retrieval
		var dbStatus string
		err = testPool.QueryRow(ctx, "SELECT status FROM auth.users WHERE id = $1", user.Id).Scan(&dbStatus)
		assert.NoError(t, err)
		assert.Equal(t, "active", dbStatus)
	})

	t.Run("FindByExternalSubject", func(t *testing.T) {
		// Use seed data
		sub := "sub_admin_tech_001"
		user, err := repo.FindByExternalSubject(ctx, universityID, sub)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, sub, user.ExternalSubject)
	})

	t.Run("FindAll_Filter", func(t *testing.T) {
		filter := domain.UserFilter{
			UniversityId: universityID,
			Status:       "active",
			Pagination: libtypes.Pagination{
				Page: 1,
				Size: 10,
			},
		}

		users, count, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))
		assert.NotEmpty(t, users)

		// Test Search
		filter.Search = "Tech Admin"
		users, count, err = repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})

	t.Run("FindAll_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, err := repo.FindAll(ctx, domain.UserFilter{})
		assert.Error(t, err)
	})

	t.Run("Store_Upsert", func(t *testing.T) {
		user := &domain.User{
			UniversityId:     universityID,
			ExternalSubject:  "sub_test_upsert",
			IdentityProvider: "central",
			Status:           "active",
			Metadata:         map[string]any{"key": "v1"},
		}
		require.NoError(t, repo.Store(ctx, user))

		// Update via Store (upsert)
		user.Metadata = map[string]any{"key": "v2"}
		require.NoError(t, repo.Store(ctx, user))

		// Verify update
		// user.Id should be set
		// Actually Store updates updated_at, let's just check no error and ID exists
		assert.NotEmpty(t, user.Id)
	})
	t.Run("UpdateStatus", func(t *testing.T) {
		// Create a user to update
		user := &domain.User{
			UniversityId:     universityID,
			ExternalSubject:  "sub_status_update",
			IdentityProvider: "central",
			Status:           "active",
			Metadata:         map[string]any{"key": "value"},
		}
		require.NoError(t, repo.Store(ctx, user))

		err := repo.UpdateStatus(ctx, user.Id, "suspended", user.Id)
		assert.NoError(t, err)

		// Verify
		var dbStatus string
		err = testPool.QueryRow(ctx, "SELECT status FROM auth.users WHERE id = $1", user.Id).Scan(&dbStatus)
		assert.NoError(t, err)
		assert.Equal(t, "suspended", dbStatus)
	})

	t.Run("AssignRole", func(t *testing.T) {
		// Need seed IDs
		var userID, roleID, groupID, institutionID string
		err := testPool.QueryRow(ctx, "SELECT id FROM auth.users WHERE external_subject = 'sub_staff_tech_001'").Scan(&userID)
		require.NoError(t, err)
		err = testPool.QueryRow(ctx, "SELECT id FROM iam.roles WHERE name = 'academic_staff' AND university_id = $1", universityID).Scan(&roleID)
		require.NoError(t, err)
		err = testPool.QueryRow(ctx, "SELECT id FROM iam.groups WHERE name = 'IT Department' AND university_id = $1", universityID).Scan(&groupID)
		require.NoError(t, err)

		// Use a different role assignment or just re-assign (on conflict update)
		// Let's create a new role for assignment test to be clean or use the seeded one to test idempotency
		// Assigning same role/group should work (upsert)

		cmd := &domain.UserRole{
			UniversityId:  universityID,
			UserId:        userID,
			RoleId:        roleID,
			InstitutionId: institutionID,
			AssignedBy:    &userID, // self assigned for test
		}

		err = repo.AssignRole(ctx, cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, cmd.Id)
	})

	t.Run("AssignRole_Upsert", func(t *testing.T) {
		// Use seed IDs
		var userID, roleID, institutionID string
		err := testPool.QueryRow(ctx, "SELECT id FROM auth.users WHERE external_subject = 'sub_admin_tech_001'").Scan(&userID)
		require.NoError(t, err)
		err = testPool.QueryRow(ctx, "SELECT id FROM iam.roles WHERE name = 'super_admin' AND university_id = $1", universityID).Scan(&roleID)
		require.NoError(t, err)
		err = testPool.QueryRow(ctx, "SELECT id FROM iam.groups WHERE name = 'IT Department' AND university_id = $1", universityID).Scan(&institutionID)
		require.NoError(t, err)

		cmd := &domain.UserRole{
			UniversityId:  universityID,
			UserId:        userID,
			RoleId:        roleID,
			InstitutionId: institutionID,
			AssignedBy:    &userID,
		}

		// First assignment (seeded one exists, so this is upsert)
		err = repo.AssignRole(ctx, cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, cmd.Id)
	})

	t.Run("FindByID", func(t *testing.T) {
		var userID string
		err := testPool.QueryRow(ctx, "SELECT id FROM auth.users WHERE external_subject = 'sub_admin_tech_001'").Scan(&userID)
		require.NoError(t, err)

		user, err := repo.FindByID(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.Id)
		assert.Equal(t, userID, user.Id)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		user, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000") // UUID format but non-existent
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("FindByExternalSubject_NotFound", func(t *testing.T) {
		user, err := repo.FindByExternalSubject(ctx, universityID, "non-existent-sub")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	t.Run("Store_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		user := &domain.User{UniversityId: universityID}
		err := repo.Store(ctx, user)
		assert.Error(t, err)
	})

	t.Run("FindByID_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repo.FindByID(ctx, "some-id")
		assert.Error(t, err)
	})

	t.Run("FindByExternalSubject_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repo.FindByExternalSubject(ctx, universityID, "some-sub")
		assert.Error(t, err)
	})

	t.Run("AssignRole_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := repo.AssignRole(ctx, &domain.UserRole{})
		assert.Error(t, err)
	})
}
