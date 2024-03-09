package uuid

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/youngfs/youngfs/pkg/kv"
)

func init() {
	uuid.EnableRandPool()
}

var (
	errUUIDDuplicated = errors.New("uuid duplicated")
)

type UUID struct {
	store kv.TransactionStore
}

func NewUUID(store kv.TransactionStore) *UUID {
	return &UUID{
		store: store,
	}
}

func (u *UUID) generateID(ctx context.Context) (uuid.UUID, error) {
	id := uuid.New()
	var err error
	check := func(id []byte) error {
		return kv.DoTransaction(u.store, ctx, func(txn kv.Transaction) error {
			_, err := txn.Get(ctx, id)
			if err != kv.ErrKeyNotFound {
				if err != nil {
					return err
				} else {
					return errUUIDDuplicated
				}
			}
			return txn.Put(ctx, id, []byte{1})
		})
	}
	for err = check([]byte(id.String())); err == errUUIDDuplicated; {
		id = uuid.New()
	}
	return id, err
}

func (u *UUID) String(ctx context.Context) (string, error) {
	id, err := u.generateID(ctx)
	if err != nil {
		return "", err
	}
	return id.String(), err
}

func (u *UUID) Bytes(ctx context.Context) ([]byte, error) {
	id, err := u.generateID(ctx)
	if err != nil {
		return nil, err
	}
	return []byte(id.String()), err
}
