package db

import (
	"context"
	"fmt"
)

// CreatePaymentRequestWithAudit creates a new payment request and creates an audit record
func (store *SQLStore) CreatePaymentRequestWithAudit(ctx context.Context, userID int32, params CreatePaymentRequestParams) (PaymentRequest, error) {
	var paymentRequest PaymentRequest
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the CreatePaymentRequest function
		paymentRequest, err = q.CreatePaymentRequest(ctx, params)
		return err
	})

	return paymentRequest, err

}

// UpdatePaymentRequestWithAudit updates a payment request record and creates an audit record
func (store *SQLStore) UpdatePaymentRequestWithAudit(ctx context.Context, userID int32, params UpdatePaymentRequestParams) (PaymentRequest, error) {
	var paymentRequest PaymentRequest
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the UpdatePaymentRequest function
		paymentRequest, err = q.UpdatePaymentRequest(ctx, params)
		return err
	})

	return paymentRequest, err
}

// DeletePaymentRequestWithAudit deletes a payment request record and creates an audit record
func (store *SQLStore) DeletePaymentRequestWithAudit(ctx context.Context, userID int32, id int32) error {
	err := store.execTx(ctx, func(q *Queries) error {
		// set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		err = q.DeletePaymentRequest(ctx, id)
		return err
	})
	return err
}

// ApprovePaymentRequestWithAudit approves a payment request record and creates an audit record
func (store *SQLStore) ApprovePaymentRequestWithAudit(ctx context.Context, userID int32, params ApprovePaymentRequestParams) (PaymentRequest, error) {
	var paymentRequest PaymentRequest
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the UpdatePaymentRequest function
		paymentRequest, err = q.ApprovePaymentRequest(ctx, params)
		return err
	})

	return paymentRequest, err
}
