package db

import (
	"context"
	"fmt"
)

type CreateQuotationTxParams struct {
	Quotation      CreateQuotationParams
	QuotationItems []CreateQuotationItemParams
}

type QuotationTxResult struct {
	Quotation      Quotation
	QuotationItems []QuotationItem
}

// CreateQuotationTxn creates a new quotation and associated quotation items in a single transaction
func (store *SQLStore) CreateQuotationTxn(ctx context.Context, userID int32, arg CreateQuotationTxParams) (QuotationTxResult, error) {
	var result QuotationTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		result.Quotation, err = q.CreateQuotation(ctx, arg.Quotation)

		if err != nil {
			return err
		}

		for _, item := range arg.QuotationItems {
			item.QuotationID = result.Quotation.ID
			quotationItem, err := q.CreateQuotationItem(ctx, item)

			if err != nil {
				return err
			}

			result.QuotationItems = append(result.QuotationItems, quotationItem)
		}

		return err
	})

	return result, err
}

type UpdateQuotationTxParams struct {
	UpdateQuotation      UpdateQuotationParams
	UpdateQuotationItems []UpdateQuotationItemParams
}

type UpdateQuotationTxResult struct {
	Quotation      Quotation
	QuotationItems []QuotationItem
}

// UpdateQuotationTxn updates a quotation and associated quotation items in a single transaction
func (store *SQLStore) UpdateQuotationTxn(ctx context.Context, userID int32, arg UpdateQuotationTxParams) (UpdateQuotationTxResult, error) {
	var result UpdateQuotationTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		result.Quotation, err = q.UpdateQuotation(ctx, arg.UpdateQuotation)
		if err != nil {
			return fmt.Errorf("error updating quotation: %w", err)
		}

		for _, item := range arg.UpdateQuotationItems {
			updatedItem, err := q.UpdateQuotationItem(ctx, item)
			if err != nil {
				return fmt.Errorf("error updating quotation item: %w", err)
			}

			result.QuotationItems = append(result.QuotationItems, updatedItem)

		}
		return err

	})

	return result, err
}
