package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomQuotationItem(t *testing.T) QuotationItem {
	quotation := createRandomQuotation(t)

	arg := CreateQuotationItemParams{
		Description: "Sample description",
		Uom:         "PCS",
		Qty:         1,
		LeadTime:    "1 week",
		ItemPrice:   100.0,
		Disc:        20.0,
		UnitPrice:   80.0,
		NetPrice:    80.0,
		Currency:    "USD",
		QuotationID: quotation.ID,
	}
	quotationItem, err := testQueries.CreateQuotationItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, quotationItem)

	return quotationItem
}

func TestQueries_CreateQuotationItem(t *testing.T) {
	createRandomQuotationItem(t)
}

func TestQueries_GetQuotationItemsByQuotationID(t *testing.T) {
	q1 := createRandomQuotationItem(t)

	q2, err := testQueries.GetQuotationItemsByQuotationID(context.Background(), q1.QuotationID)
	require.NoError(t, err)
	require.NotEmpty(t, q2)
	require.Equal(t, q1.ID, q2[0].ID)
	require.Equal(t, q1.Description, q2[0].Description)
	require.Equal(t, q1.Uom, q2[0].Uom)

}
