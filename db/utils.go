package db

import (
	"context"
	"database/sql"
	"fmt"
)

func CountQuery(db *sql.DB, ctx context.Context, table string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)

	err := db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func DeleteQuery(db *sql.DB, ctx context.Context, table string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE Id = $1", table)

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
