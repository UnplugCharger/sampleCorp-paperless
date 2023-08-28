package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomQuotation(t *testing.T) Quotation {
	status := "sent"
	company := createRandomCompany(t)
	sig := createRandomSignatory(t)
	arg := CreateQuotationParams{
		Attn:           "Attn",
		CompanyID:      int32(company.ID),
		Site:           "Site",
		Validity:       1,
		Warranty:       1,
		PaymentTerms:   "PaymentTerms",
		DeliveryTerms:  "DeliveryTerms",
		SignatoryID:    int32(sig.ID),
		Status:         "PENDING",
		SentOrReceived: &status,
	}
	quotation, err := testQueries.CreateQuotation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, quotation)

	return quotation
}

func TestQueries_CreateQuotation(t *testing.T) {
	createRandomQuotation(t)
}
