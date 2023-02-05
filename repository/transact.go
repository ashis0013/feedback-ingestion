package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func transactExcec(db *sqlx.DB, txs func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Fatal("failed to rollback")
		}
	}()

	if err != nil {
		log.Fatal("Cannot begin db transaction")
		return err
	}

	if txs != nil {
		err = txs(tx)
	}
	if err != nil {
		log.Fatal("Failed to execute query")
		return err
	}
	return nil
}

func transactQuery[T any](db *sqlx.DB, records *[]*T, txs func(*sqlx.Tx) (*sqlx.Rows, error)) error {
	tx, err := db.Beginx()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Fatal("failed to rollback")
		}
	}()
	if err != nil {
		log.Fatal("Cannot begin db transaction")
		return err
	}

	if txs == nil {
		return err
	}

	rows, err := txs(tx)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record T
		err := rows.StructScan(&record)
		if err != nil {
			log.Fatal("Error while reading single record")
			return err
		}
		*records = append(*records, &record)
	}
	return nil
}
