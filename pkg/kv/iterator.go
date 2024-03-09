package kv

type Iterator interface {
	Seek(key []byte)
	Valid() bool
	Next()
	Key() []byte
	Value() []byte
	Close()
}
