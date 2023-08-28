package db

import (
	"context"
	"fmt"
)

func (store *SQLStore) CreatePettyCashWithAudit(ctx context.Context, userID int32, params CreatePettyCashParams) (PettyCash, error) {
	var pettyCash PettyCash
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return fmt.Errorf("error setting audit.current_user_id: %v", err)
		}

		// Call the CreatePettyCash function
		pettyCash, err = q.CreatePettyCash(ctx, params)
		if err != nil {
			return fmt.Errorf("error calling CreatePettyCash: %v", err)
		}

		return nil
	})

	return pettyCash, err
}

// UpdatePettyCashWithAudit updates a petty cash record and creates an audit record
func (store *SQLStore) UpdatePettyCashWithAudit(ctx context.Context, userID int32, params UpdatePettyCashParams) (PettyCash, error) {
	var pettyCash PettyCash
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the UpdatePettyCash function
		pettyCash, err = q.UpdatePettyCash(ctx, params)
		return err
	})

	return pettyCash, err
}

// DeletePettyCashWithAudit deletes a petty cash record and creates an audit record
func (store *SQLStore) DeletePettyCashWithAudit(ctx context.Context, userID int32, id int32) error {

	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the DeletePettyCash function
		err = q.DeletePettyCash(ctx, id)
		return err
	})

	return err
}

// ApprovePettyCashWithAudit ApprovePettyCashRequestWithAudit approves a petty cash request record and creates an audit record
func (store *SQLStore) ApprovePettyCashWithAudit(ctx context.Context, userID int32, params ApprovePettyCashParams) (PettyCash, error) {
	var pettyCash PettyCash
	err := store.execTx(ctx, func(q *Queries) error {
		// Set the audit.current_user_id variable
		_, err := q.db.Exec(ctx, fmt.Sprintf("SET LOCAL audit.current_user_id = %d", userID))
		if err != nil {
			return err
		}

		// Call the ApprovePettyCash function
		pettyCash, err = q.ApprovePettyCash(ctx, params)
		return err
	})

	return pettyCash, err

}
