package kv

import (
	"github.com/dgraph-io/badger/v2"
)

type BadgerKV struct {
	db *badger.DB
}

func (p *BadgerKV) Get(key []byte, opValue func([]byte) error) error {
	return p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return opValue(val)
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *BadgerKV) Put(key []byte, value []byte) error {
	txn := p.db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (p *BadgerKV) Del(key []byte) error {
	txn := p.db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (p *BadgerKV) Close() error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return err
		}
	}
	p.db = nil
	return nil
}

func OpenBadgerKV(path string) (KV, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, err
	}
	return &BadgerKV{
		db: db,
	}, nil
}
