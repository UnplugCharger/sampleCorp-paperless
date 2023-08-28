package db

import (
	"context"
	"testing"

	"github.com/qwetu_petro/backend/utils"

	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) Role {
	name := utils.RandomString(10)
	description := utils.RandomString(10)
	role, err := testQueries.CreateRole(context.Background(), CreateRoleParams{
		Name:        name,
		Description: &description,
	})
	require.NoError(t, err)
	require.NotEmpty(t, role)
	return role
}

func TestQueries_CreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestQueries_GetRole(t *testing.T) {
	role := createRandomRole(t)
	role2, err := testQueries.GetRole(context.Background(), role.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)
	require.Equal(t, role.ID, role2.ID)
	require.Equal(t, role.Name, role2.Name)
	require.Equal(t, role.Description, role2.Description)
}

func TestQueries_UpdateRole(t *testing.T) {
	role := createRandomRole(t)
	name := utils.RandomString(10)
	description := utils.RandomString(10)

	role2, err := testQueries.UpdateRole(context.Background(), UpdateRoleParams{
		ID:          role.ID,
		Name:        name,
		Description: &description,
	})
	require.NoError(t, err)
	require.NotEmpty(t, role2)
	require.Equal(t, role.ID, role2.ID)
	require.Equal(t, name, role2.Name)
	require.Equal(t, description, *role2.Description)

}

func TestQueries_DeleteRole(t *testing.T) {
	role := createRandomRole(t)
	err := testQueries.DeleteRole(context.Background(), role.ID)
	require.NoError(t, err)
	role2, err := testQueries.GetRole(context.Background(), role.ID)
	require.Error(t, err)
	require.Empty(t, role2)
}

func TestQueries_ListRoles(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRole(t)
	}
	arg := ListRolesParams{
		Limit:  5,
		Offset: 5,
	}
	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, 5)
	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}
