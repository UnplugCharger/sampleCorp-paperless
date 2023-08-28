package db

import (
	"context"
	"testing"

	"github.com/qwetu_petro/backend/utils"

	"github.com/stretchr/testify/require"
)

// TO DO  make random signatory  reference  users table
func createRandomSignatory(t *testing.T) Signatory {
	arg := CreateSignatoryParams{
		Name:  utils.RandomString(10),
		Title: utils.RandomString(10),
	}
	signatory, err := testQueries.CreateSignatory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, signatory)
	return signatory
}

func TestQueries_CreateSignatory(t *testing.T) {
	createRandomSignatory(t)
}

func TestQueries_DeleteSignatoryByName(t *testing.T) {
	signatory := createRandomSignatory(t)
	err := testQueries.DeleteSignatoryByName(context.Background(), signatory.Name)
	require.NoError(t, err)

}
