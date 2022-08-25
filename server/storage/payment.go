package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
)

type Payment struct {
	PreimageHash string `db:"PreimageHash"`
	Preimage     string `db:"Preimage"`
	Invoice      string `db:"Invoice"`
	Amount       int    `db:"Amount"`
}

type PaymentRequest struct {
	Amount int `json:"amount"`
}

const createPaymentTable = `
CREATE TABLE IF NOT EXISTS payment (
	PreimageHash TEXT PRIMARY KEY,
	Preimage TEXT,
	Invoice TEXT,
	Amount TEXT
)
`

const insertPayment = `
INSERT INTO payment (
    PreimageHash,
	Preimage,
	Invoice,
	Amount
)
VALUES (
    :PreimageHash,
	:Preimage,
	:Invoice,
	:Amount
)
`

const getPayments = `
SELECT *
FROM payment`

const getPayment = `
SELECT *
FROM payment
WHERE PreimageHash = ?
LIMIT 1`

func (db *DB) SavePayment(payment *Payment) error {
	log.Debug("inserting payment{", payment.PreimageHash, "}", ": ", payment)
	query, args, _ := sqlx.Named(insertPayment, payment)
	result, err := db.db.Exec(query, args...)
	if err != nil {
		log.Error("Error inserting payment", err)
		return err
	}
	log.Debug("success inserting payment", result)
	return nil
}

func (db *DB) GetPayments() ([]*Payment, error) {
	log.Debug("retrieving payments... ")
	var payments []*Payment
	rows, err := db.db.Queryx(getPayments)
	if err == sql.ErrNoRows {
		return payments, nil
	}

	for rows.Next() {
		var payment Payment

		if err := rows.StructScan(&payment); err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}
	return payments, nil
}

func (db *DB) GetPayment(preimageHash string) (result Payment, err error) {
	log.Debugf("retrieving payment: %s ", preimageHash)
	rows, err := db.db.Queryx(getPayment, preimageHash)
	if err != nil {
		log.Fatalf("Error retrieving payment: %v", err)
		return
	}
	if rows.Next() {
		if err = rows.StructScan(&result); err != nil {
			log.Fatalf("Error mapping payment: %v", err)
			return
		}
	}
	return

}
