package integrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	libtypes "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
	rolesRepo "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/repository/postgresql"
)

func TestRolesRepository(t *testing.T) {
	if testPool == nil {
		t.Skip("Skipping integration test: testPool is nil")
	}
	ctx := context.Background()
	repo := rolesRepo.NewRepository(testPool)

	var instID string
	err := testPool.QueryRow(ctx, "SELECT id FROM auth.institutions WHERE code = 'TECH-UNI'").Scan(&instID)
	require.NoError(t, err)

	var userID string
	err = testPool.QueryRow(ctx, "SELECT id FROM auth.users WHERE external_subject = 'sub_admin_tech_001'").Scan(&userID)
	require.NoError(t, err)

	t.Run("Roles_CRUD", func(t *testing.T) {
		// Create
		role := &domain.Role{
			UniversityId: instID,
			Name:         "custom_role",
			Description:  "Custom Role Desc",
			IsActive:     true,
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}

		err := repo.Store(ctx, role) // Note: Store might not be returning ID on struct if not implemented, check repo.
		// Checking repo store logic... usually it does set ID.
		// If not, we query by name.
		assert.NoError(t, err)

		// If ID is empty after store (check implementation), query it.
		if role.Id == "" {
			// Try to find it
			found, err := repo.FindByName(ctx, instID, "custom_role")
			require.NoError(t, err)
			role.Id = found.Id
		}

		assert.NotEmpty(t, role.Id)

		// Read
		found, err := repo.FindByID(ctx, role.Id)
		assert.NoError(t, err)
		assert.Equal(t, "custom_role", found.Name)

		// Update
		role.Description = "Updated Desc"
		err = repo.Update(ctx, role)
		assert.NoError(t, err)

		found, err = repo.FindByID(ctx, role.Id)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Desc", found.Description)

		// Delete
		err = repo.Delete(ctx, instID, role.Id)
		assert.NoError(t, err)

		_, err = repo.FindByID(ctx, role.Id)
		assert.Error(t, err) // Should be not found error
	})

	t.Run("FindAll", func(t *testing.T) {
		filter := domain.RoleFilter{
			UniversityId: instID,
			Pagination:   libtypes.Pagination{Page: 1, Size: 10},
		}
		roles, count, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))
		assert.NotEmpty(t, roles)
	})

	t.Run("Permissions_Management", func(t *testing.T) {
		// Create a role for permissions test
		role := &domain.Role{
			UniversityId: instID,
			Name:         "perm_test_role",
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}
		require.NoError(t, repo.Store(ctx, role))
		// Ensure ID
		if role.Id == "" {
			r, _ := repo.FindByName(ctx, instID, "perm_test_role")
			role.Id = r.Id
		}

		// Add Permissions
		perms := []string{"users.manage.all.view", "roles.manage.all.view"}
		err := repo.AddPermissions(ctx, role.Id, perms)
		assert.NoError(t, err)

		// Get Permissions
		fetchedPerms, err := repo.GetPermissions(ctx, role.Id)
		assert.NoError(t, err)
		assert.ElementsMatch(t, perms, fetchedPerms)

		// Remove Permissions
		err = repo.RemovePermissions(ctx, role.Id)
		assert.NoError(t, err)

		fetchedPerms, err = repo.GetPermissions(ctx, role.Id)
		assert.NoError(t, err)
		assert.Empty(t, fetchedPerms)
	})

	t.Run("FindAllPermissions", func(t *testing.T) {
		filter := domain.PermissionFilter{
			UniversityId: instID,
			Pagination:   libtypes.Pagination{Page: 1, Size: 10},
		}
		perms, err := repo.FindAllPermissions(ctx, filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, perms)
	})

	t.Run("FindAll_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, err := repo.FindAll(ctx, domain.RoleFilter{})
		assert.Error(t, err)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		role, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Nil(t, role)
	})

	t.Run("AddPermissions_Edges", func(t *testing.T) {
		// Empty list
		err := repo.AddPermissions(ctx, "some-role", []string{})
		assert.NoError(t, err)

		// Invalid permissions (not found)
		// Need a real role ID though to avoid FK error if it inserts?
		// Actually my logic returns early if no permission IDs found.
		// So role ID validity doesn't matter for this specific path check (len=0).
		err = repo.AddPermissions(ctx, "some-role", []string{"invalid.perm"})
		assert.NoError(t, err)
	})
	t.Run("Delete_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := repo.Delete(ctx, instID, "some-id")
		assert.Error(t, err)
	})

	t.Run("FindByID_ContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repo.FindByID(ctx, "some-id")
		assert.Error(t, err)
	})
}
