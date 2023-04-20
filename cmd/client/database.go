package main

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type Database struct {
	dbClient *badger.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) Open() error {
	badgerDb, err := badger.Open(badger.DefaultOptions(defaultBadgerLocation))
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	db.dbClient = badgerDb
	return nil
}

func (db *Database) Put(key, value []byte) error {
	return db.dbClient.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)
		return txn.SetEntry(entry)
	})
}

func (db *Database) Get(key []byte) (result []byte, err error) {
	err = db.dbClient.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			result = val
			return nil
		})
	})

	return result, err
}

func (db *Database) Close() error {
	return db.dbClient.Close()
}
