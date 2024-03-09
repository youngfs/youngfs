package tikv

import (
	"context"
	"github.com/youngfs/youngfs/pkg/kv"

	tikverr "github.com/tikv/client-go/v2/error"
	"github.com/tikv/client-go/v2/txnkv/transaction"
)

func (s *Store) Put(ctx context.Context, key []byte, val []byte) error {
	return s.update(ctx, func(txn *transaction.KVTxn) error {
		return txn.Set(key, val)
	})
}

func (s *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	var val []byte
	return val, s.view(func(txn *transaction.KVTxn) error {
		var err error
		val, err = txn.Get(ctx, key)
		if tikverr.IsErrNotFound(err) {
			return kv.ErrKeyNotFound
		}
		return err
	})
}

func (s *Store) Delete(ctx context.Context, key []byte) error {
	return s.update(ctx, func(txn *transaction.KVTxn) error {
		return txn.Delete(key)
	})
}
