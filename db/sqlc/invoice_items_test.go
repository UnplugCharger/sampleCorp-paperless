package db

import (
	"context"
	"testing"

	"github.com/qwetu_petro/backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomInvoiceItem(t *testing.T) InvoiceItem {
	invoice := createRandomInvoice(t)
	arg := CreateInvoiceItemParams{
		Description: utils.RandomString(10),
		Uom:         utils.RandomString(10),
		Qty:         int32(utils.RandomInt()),
		UnitPrice:   1000, // Generate a random number between 0 and 10000
		NetPrice:    100,  // Generate a random number between 0 and 10000
		Currency:    utils.RandomString(10),
		InvoiceID:   invoice.ID,
	}
	invoiceItem, err := testQueries.CreateInvoiceItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoiceItem)
	return invoiceItem
}

func TestQueries_CreateInvoiceItem(t *testing.T) {
	createRandomInvoiceItem(t)
}

func TestQueries_GetInvoiceItemsByInvoiceID(t *testing.T) {
	invoiceItem1 := createRandomInvoiceItem(t)
	invoiceItem2, err := testQueries.GetInvoiceItemsByInvoiceID(context.Background(), invoiceItem1.InvoiceID)
	require.NoError(t, err)
	require.NotEmpty(t, invoiceItem2)
}
