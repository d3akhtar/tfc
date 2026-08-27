package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxFunc func(tx *sql.Tx) error

func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transcation: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Transcation failed: %w / Rollback failed: %w", err, rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Transcation commit failed: %w", err)
	}

	return nil
}
