package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	driver = "sqlite3"
)

type DBConnector interface {
	CheckConnection() error
	Close() error
	SavePayment(payment *Payment) error
	GetPayments() (result []*Payment, err error)
	GetPayment(preimageHash string) (result Payment, err error)
	SaveConfig(config *Config) error
	GetConfigs() (result []*Config, err error)
	GetConfig(key string) (result Config, err error)
	SetDefaultConfig() error
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

	if _, err := db.Exec(createConfigTable); err != nil {
		return nil, err
	}

	result := &DB{db}
	return result, nil
}

func (db *DB) Close() error {
	log.Debug("closing connection to DB")
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}
