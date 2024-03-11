package needle

import (
	"context"
	"github.com/youngfs/youngfs/pkg/kv"
	"sync/atomic"
)

type KvNeedleStore struct {
	kv   kv.Store
	size uint64
	ctx  context.Context
}

func NewKvStore(kv kv.Store) (*KvNeedleStore, error) {
	it, err := kv.NewIterator()
	if err != nil {
		return nil, err
	}
	defer it.Close()
	size := uint64(0)
	for ; it.Valid(); it.Next() {
		size++
	}
	return &KvNeedleStore{kv: kv, size: size, ctx: context.Background()}, nil
}

func (s *KvNeedleStore) Put(n *Needle) error {
	err := s.kv.Put(s.ctx, n.Id.Bytes(), n.ToBytes())
	if err != nil {
		return err
	}
	atomic.AddUint64(&s.size, 1)
	return nil
}

func (s *KvNeedleStore) Get(id Id) (*Needle, error) {
	b, err := s.kv.Get(s.ctx, id.Bytes())
	if err != nil {
		return nil, err
	}
	return FromBytes(b)
}

func (s *KvNeedleStore) Delete(id Id) error {
	err := s.kv.Delete(s.ctx, id.Bytes())
	if err != nil {
		return err
	}
	atomic.AddUint64(&s.size, -1)
	return nil
}

func (s *KvNeedleStore) Size() uint64 {
	return s.size
}

func (s *KvNeedleStore) Close() error {
	return s.kv.Close()
}

func KvNeedleStoreCreator(kvCreator func(path string) (kv.Store, error)) func(path string) (Store, error) {
	return func(path string) (Store, error) {
		kvStore, err := kvCreator(path)
		if err != nil {
			return nil, err
		}
		return NewKvStore(kvStore)
	}
}
