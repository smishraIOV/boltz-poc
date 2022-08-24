package storage

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const (
	driver = "sqlite"
)

type DBConnector interface {
	CheckConnection() error
	Close() error
	SavePayment(payment *Payment) error
	GetPayments() (result []*Payment, err error)
}

type DB struct {
	db *sqlx.DB
}

func (db *DB) CheckConnection() error {
	return db.db.Ping()
}

func Connect(dbPath string) (*DB, error) {
	log.Debug("connecting to DB: ", dbPath)
	db, err := sqlx.Connect(driver, dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createPaymentTable); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	log.Debug("closing connection to DB")
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}
