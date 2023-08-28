package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store provides all functions  to run database queries and transactions
type Store interface {
	Querier
	CreateInvoiceTxn(ctx context.Context, userID int32, arg1 CreateInvoiceParams, arg2 []CreateInvoiceItemParams) (InvoiceTxResult, error)
	CreatePurchaseOrderTxn(ctx context.Context, arg CreatePurchaseOrderTxParams, userID int32) (PurchaseOrderTxResult, error)
	UpdatePurchaseOrderTxn(ctx context.Context, arg UpdatePurchaseOrderTxParams, userID int32) (PurchaseOrderTxResult, error)
	ApprovePurchaseOrderTxn(ctx context.Context, arg ApprovePurchaseOrderParams, userID int32) (PurchaseOrderTxResult, error)
	CreateQuotationTxn(ctx context.Context, userID int32, arg CreateQuotationTxParams) (QuotationTxResult, error)
	UpdateQuotationTxn(ctx context.Context, userID int32, arg UpdateQuotationTxParams) (UpdateQuotationTxResult, error)
	CreatePettyCashWithAudit(ctx context.Context, userID int32, params CreatePettyCashParams) (PettyCash, error)
	UpdatePettyCashWithAudit(ctx context.Context, userID int32, params UpdatePettyCashParams) (PettyCash, error)
	DeletePettyCashWithAudit(ctx context.Context, userID int32, id int32) error
	CreatePaymentRequestWithAudit(ctx context.Context, userID int32, params CreatePaymentRequestParams) (PaymentRequest, error)
	UpdatePaymentRequestWithAudit(ctx context.Context, userID int32, params UpdatePaymentRequestParams) (PaymentRequest, error)
	ApprovePaymentRequestWithAudit(ctx context.Context, userID int32, params ApprovePaymentRequestParams) (PaymentRequest, error)
	ApprovePettyCashWithAudit(ctx context.Context, userID int32, params ApprovePettyCashParams) (PettyCash, error)
	DeletePaymentRequestWithAudit(ctx context.Context, userID int32, id int32) error
}

// SQLStore / SQLStore provides all functions to run database queries and transactions
type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

// NewStore creates a new Store with connection pooling
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(pool),
		pool:    pool,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	conn, err := store.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
