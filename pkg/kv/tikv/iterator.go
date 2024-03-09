package tikv

import (
	"github.com/tikv/client-go/v2/tikv"
	"github.com/tikv/client-go/v2/txnkv/transaction"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/util"
)

type Iterator struct {
	txn    *transaction.KVTxn
	iter   tikv.Iterator
	prefix []byte
	end    []byte
}

func (s *Store) NewIterator(opts ...kv.IteratorOption) (kv.Iterator, error) {
	cfg := &kv.IteratorConfig{}
	kv.ApplyConfig(cfg, opts...)

	txn, err := s.getTxn()
	if err != nil {
		return nil, err
	}
	var end []byte
	if cfg.Prefix != nil {
		end = util.PrefixEnd(cfg.Prefix)
	}

	return &Iterator{
		txn:    txn,
		prefix: cfg.Prefix,
		end:    end,
	}, nil
}

func (it *Iterator) Seek(key []byte) {
	if it.iter != nil {
		it.iter.Close()
	}
	iter, err := it.txn.Iter(key, it.end)
	if err != nil {
		panic(err)
	}
	it.iter = iter
}

func (it *Iterator) Valid() bool {
	if it.iter == nil {
		return false
	}
	return it.iter.Valid()
}

func (it *Iterator) Next() {
	if it.iter != nil {
		err := it.iter.Next()
		if err != nil {
			panic(err)
		}
	}
}

func (it *Iterator) Key() []byte {
	if it.iter == nil {
		return nil
	}
	return it.iter.Key()
}

func (it *Iterator) Value() []byte {
	if it.iter == nil {
		return nil
	}
	return it.iter.Value()
}

func (it *Iterator) Close() {
	it.iter.Close()
}
