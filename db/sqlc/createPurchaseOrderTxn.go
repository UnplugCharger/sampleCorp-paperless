package db

import (
	"context"
	"fmt"
)

type CreatePurchaseOrderTxParams struct {
	PurchaseOrder      CreatePurchaseOrderParams
	PurchaseOrderItems []CreatePurchaseOrderItemParams
}

type PurchaseOrderTxResult struct {
	PurchaseOrder      PurchaseOrder
	PurchaseOrderItems []PurchaseOrderItem
}

// CreatePurchaseOrderTxn creates a new purchase order and associated purchase order items in a single transaction

func (store *SQLStore) CreatePurchaseOrderTxn(ctx context.Context, arg CreatePurchaseOrderTxParams, userID int32) (PurchaseOrderTxResult, error) {
	var result PurchaseOrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		result.PurchaseOrder, err = q.CreatePurchaseOrder(ctx, arg.PurchaseOrder)
		if err != nil {
			return err
		}

		for _, item := range arg.PurchaseOrderItems {
			id := result.PurchaseOrder.ID
			item.PurchaseOrderID = id
			purchaseOrderItem, err := q.CreatePurchaseOrderItem(ctx, item)
			if err != nil {
				return err
			}
			result.PurchaseOrderItems = append(result.PurchaseOrderItems, purchaseOrderItem)
		}
		return err

	})

	return result, err
}

type UpdatePurchaseOrderTxParams struct {
	PurchaseOrder      UpdatePurchaseOrderParams
	PurchaseOrderItems []UpdatePurchaseOrderItemParams
}

// UpdatePurchaseOrderTxn creates a new purchase order and associated purchase order items in a single transaction
func (store *SQLStore) UpdatePurchaseOrderTxn(ctx context.Context, arg UpdatePurchaseOrderTxParams, userID int32) (PurchaseOrderTxResult, error) {
	var result PurchaseOrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		result.PurchaseOrder, err = q.UpdatePurchaseOrder(ctx, arg.PurchaseOrder)
		if err != nil {
			fmt.Println("Error updating purchase order")
			return err
		}

		for _, item := range arg.PurchaseOrderItems {
			id := result.PurchaseOrder.ID
			item.PurchaseOrderID = id
			purchaseOrderItem, err := q.UpdatePurchaseOrderItem(ctx, item)
			if err != nil {
				fmt.Println("Error updating purchase order item")
				return err
			}
			result.PurchaseOrderItems = append(result.PurchaseOrderItems, purchaseOrderItem)
		}
		return err

	})

	return result, err
}

// ApprovePurchaseOrderTxn creates a new purchase order and associated purchase order items in a single transaction
func (store *SQLStore) ApprovePurchaseOrderTxn(ctx context.Context, arg ApprovePurchaseOrderParams, userID int32) (PurchaseOrderTxResult, error) {
	var result PurchaseOrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Set the audit.current_user_id variable
		_, err = q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		result.PurchaseOrder, err = q.ApprovePurchaseOrder(ctx, arg)
		return err
	})
	return result, err

}
