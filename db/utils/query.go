package utils

import (
	"context"
	"database/sql"
	"fmt"

	database "github.com/d3akhtar/tfc/db"
)

func CountQuery(db *sql.DB, ctx context.Context, table string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)

	err := db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func ExecStmtUpdate(stmt *sql.Stmt, ctx context.Context, args ...any) error {
	return checkResultFromUpdateCall(func() (sql.Result, error) { return stmt.ExecContext(ctx, args...) })
}

func ExecQueryUpdate(query string, db *sql.DB, ctx context.Context, args ...any) error {
	return checkResultFromUpdateCall(func() (sql.Result, error) { return db.ExecContext(ctx, query, args...) })
}

func checkResultFromUpdateCall(update func() (sql.Result, error)) error {
	result, err := update()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func DeleteQuery(db *sql.DB, ctx context.Context, table string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE Id = $1", table)
	return ExecQueryUpdate(query, db, ctx, id)
}
