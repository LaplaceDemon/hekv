package kv

import (
	"github.com/cockroachdb/pebble"
)

type PebbleKV struct {
	db   *pebble.DB
	opts *pebble.Options
	ro   *pebble.IterOptions
	wo   *pebble.WriteOptions
}

func (p *PebbleKV) Get(key []byte, op func([]byte) error) (err error) {
	val, closer, err := p.db.Get(key)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	defer func() {
		if closer != nil {
			if cerr := closer.Close(); err == nil {
				err = cerr
			}
		}
	}()
	return op(val)
}

func (p *PebbleKV) Put(key []byte, value []byte) error {
	return p.db.Set(key, value, p.wo)
}

func (p *PebbleKV) Del(key []byte) error {
	return p.db.Delete(key, p.wo)
}

func (p *PebbleKV) Close() error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return err
		}
	}
	p.db = nil
	return nil
}

func OpenPebbleKV(dir string) (KV, error) {
	opts := &pebble.Options{}

	db, err := pebble.Open(dir, opts)
	if err != nil {
		return nil, err
	}

	ro := &pebble.IterOptions{}
	wo := &pebble.WriteOptions{Sync: true}
	return &PebbleKV{
		db:   db,
		ro:   ro,
		wo:   wo,
		opts: opts,
	}, nil
}
