package dbutil

import (
	"context"
	"errors"
	"regexp"

	"github.com/bernardinorafael/gogem/fault"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// VerifyDuplicatedKey checks if the error is a unique constraint violation
//
// This function analyzes PostgreSQL errors to identify unique constraint violations
// (error code 23505). When a violation is detected, it extracts the field name
// that caused the conflict and returns a custom error.
//
// Parameters:
//   - target: the error to be analyzed (must be of type *pq.Error)
//
// Returns:
//   - string: name of the field that caused the violation (empty if not a violation)
//   - error: custom error with Database tag if violation, nil otherwise
func VerifyDuplicatedKey(target error) (string, error) {
	var pqErr *pq.Error
	if errors.As(target, &pqErr) && pqErr.Code == "23505" {
		// Compile the regular expression to find the pattern "Key (field)="
		re := regexp.MustCompile(`(?i)Key\s*\(\s*(.*?)\s*\)\s*=`)
		matches := re.FindStringSubmatch(pqErr.Detail)
		field := matches[1]

		return field, fault.New("duplicated constraint key", fault.WithTag(fault.DB))
	}

	return "", nil
}

// VerifyForeignKeyViolation checks if the error is a foreign key constraint violation
//
// This function analyzes PostgreSQL errors to identify foreign key constraint violations
// (error code 23503). When a violation is detected, it extracts the field name and
// constraint name that caused the conflict and returns a custom error.
//
// Parameters:
//   - target: the error to be analyzed (must be of type *pq.Error)
//
// Returns:
//   - string: name of the field that caused the violation (empty if not a violation)
//   - error: custom error with Database tag if violation, nil otherwise
func VerifyForeignKeyViolation(target error) (string, error) {
	var pqErr *pq.Error
	if errors.As(target, &pqErr) && pqErr.Code == "23503" {
		// Format: "Key (field_name)=(value) is not present in table \"table_name\""
		re := regexp.MustCompile(`(?i)Key\s*\(\s*(.*?)\s*\)\s*=`)
		matches := re.FindStringSubmatch(pqErr.Detail)
		field := matches[1]

		return field, fault.New("foreign key constraint violation", fault.WithTag(fault.DB))
	}

	return "", nil
}

func ExecTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fault.New("failed to begin transaction", fault.WithTag(fault.TX), fault.WithErr(err))
	}

	if err = fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fault.New("failed to rollback transaction", fault.WithTag(fault.TX), fault.WithErr(rollbackErr))
		}
		return fault.New("transaction failed", fault.WithTag(fault.TX), fault.WithErr(err))
	}

	if err := tx.Commit(); err != nil {
		return fault.New("failed to commit transaction", fault.WithTag(fault.TX), fault.WithErr(err))
	}

	return nil
}
