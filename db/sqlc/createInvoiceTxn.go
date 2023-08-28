package db

import (
	"context"
	"fmt"
)

type CreateInvoiceTxParams struct {
	Invoice      CreateInvoiceParams
	InvoiceItems []CreateInvoiceItemParams
}

type InvoiceTxResult struct {
	Invoice      Invoice
	InvoiceItems []InvoiceItem
}

// CreateInvoiceTxn creates a new invoice and associated invoice items in a single transaction
func (store *SQLStore) CreateInvoiceTxn(ctx context.Context, userID int32, invoiceParams CreateInvoiceParams, invoiceItemsParams []CreateInvoiceItemParams) (InvoiceTxResult, error) {
	var result InvoiceTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		invoice, err := q.CreateInvoice(ctx, invoiceParams)
		if err != nil {
			return err
		}

		fmt.Println("invoice: ", invoice)

		result.Invoice = invoice
		// loop  through the list of params and create invoice items
		for _, item := range invoiceItemsParams {

			item.InvoiceID = invoice.ID
			invoiceItem, err := q.CreateInvoiceItem(ctx, item)
			if err != nil {
				return err
			}
			result.InvoiceItems = append(result.InvoiceItems, invoiceItem)
		}

		return err
	})

	return result, err
}
