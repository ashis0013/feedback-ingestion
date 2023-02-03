package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func transactExcec(db *sqlx.DB, txs func(*sqlx.Tx) error) {
	tx, err := db.Beginx()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Fatal("failed to rollback")
		}
	}()

	if err != nil {
		log.Fatal("Cannot begin db transaction")
		return
	}

	if txs != nil {
		err = txs(tx)
	}
	if err != nil {
		log.Fatal("Failed to execute query")
	}
}

func transactQuery[T any](db *sqlx.DB, records *[]*T, txs func(*sqlx.Tx) (*sqlx.Rows, error)) {
	tx, err := db.Beginx()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Fatal("failed to rollback")
		}
	}()
	if err != nil {
		log.Fatal("Cannot begin db transaction")
		return
	}

	if txs == nil {
		return
	}

	rows, err := txs(tx)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var record T
		err := rows.StructScan(&record)
		if err != nil {
			log.Fatal("Error while reading single record")
		}
		*records = append(*records, &record)
	}
}
