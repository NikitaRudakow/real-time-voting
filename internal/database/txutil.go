package database

import (
	"database/sql"
	"log"
)

func RollbackTx(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Printf("failed to rollback transaction: %v", err)
	}
}
