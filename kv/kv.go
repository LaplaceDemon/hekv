package kv

type KV interface {
	Get(key []byte, opValue func([]byte) error) error
	Put(key []byte, value []byte) error
	Del(key []byte) error
	Close() error
}
