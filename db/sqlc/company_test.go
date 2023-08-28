package db

import (
	"context"
	"testing"

	"github.com/qwetu_petro/backend/utils"
)

func createRandomCompany(t *testing.T) Company {
	addr := utils.RandomString(100)
	arg := CreateCompanyParams{
		Name:     utils.RandomString(10),
		Initials: utils.RandomString(3),
		Address:  &addr,
	}
	company, err := testQueries.CreateCompany(context.Background(), arg)
	if err != nil {
		t.Fatal(err)
	}
	return company
}

func TestQueries_CreateCompany(t *testing.T) {
	createRandomCompany(t)
}

func TestQueries_DeleteCompanyByName(t *testing.T) {
	company := createRandomCompany(t)
	err := testQueries.DeleteCompanyByName(context.Background(), company.Name)
	if err != nil {
		t.Fatal(err)
	}
}
