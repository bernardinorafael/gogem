// Package dbutil provides PostgreSQL helper functions for constraint violation
// detection, a JSONB type for structured column data, and a transaction wrapper.
//
// Detecting unique constraint violations:
//
//	field, err := dbutil.VerifyDuplicatedKey(dbErr)
//	if err != nil {
//	    // field = "email", err = fault with Database tag
//	    return fault.NewConflict(fmt.Sprintf("%s already exists", field))
//	}
//
// Detecting foreign key violations:
//
//	field, err := dbutil.VerifyForeignKeyViolation(dbErr)
//	if err != nil {
//	    // field = "organization_id", err = fault with Database tag
//	    return fault.NewBadRequest(fmt.Sprintf("invalid %s reference", field))
//	}
//
// Transaction wrapper with automatic rollback on error:
//
//	err := dbutil.ExecTx(ctx, db, func(tx *sqlx.Tx) error {
//	    _, err := tx.ExecContext(ctx, "INSERT INTO users ...")
//	    if err != nil {
//	        return err // triggers rollback
//	    }
//	    _, err = tx.ExecContext(ctx, "INSERT INTO profiles ...")
//	    return err
//	})
//
// JSONB type for PostgreSQL jsonb columns:
//
//	type User struct {
//	    ID       string      `db:"id"`
//	    Metadata dbutil.JSONB `db:"metadata"` // map[string]string stored as jsonb
//	}
package dbutil
