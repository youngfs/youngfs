package uuid

import (
	"context"
	"github.com/google/uuid"
	"github.com/youngfs/youngfs/pkg/kv"
	"sync"
	"time"
)

type BatchUUID struct {
	store      kv.TransactionStore
	batchSize  uint64
	batchUUIDS []uuid.UUID
	lock       sync.Mutex
	size       uint64
	retry      int
	ttl        *time.Duration
	lastTime   time.Time
}

func NewBatchUUID(store kv.TransactionStore, batchSize uint64, ttl *time.Duration) *BatchUUID {
	return &BatchUUID{
		store:      store,
		batchSize:  batchSize,
		batchUUIDS: make([]uuid.UUID, batchSize),
		lock:       sync.Mutex{},
		size:       0,
		retry:      3,
		ttl:        ttl,
	}
}

func (u *BatchUUID) generateIDs(ctx context.Context) error {
	var err error
	for i := 0; i < u.retry; i++ {
		err = kv.DoTransaction(u.store, ctx, func(txn kv.Transaction) error {
			for i := uint64(0); i < u.batchSize; i++ {
				id := uuid.New()
				key := []byte(id.String())
				_, err := txn.Get(ctx, key)
				if err != kv.ErrKeyNotFound {
					if err != nil {
						return err
					} else {
						i--
						continue
					}
				}
				err = txn.Put(ctx, key, []byte{1})
				if err != nil {
					return err
				}
				u.batchUUIDS[i] = id
			}
			return nil
		})
		if err == nil {
			u.size = u.batchSize
			u.lastTime = time.Now()
			return nil
		}
	}
	return err
}

func (u *BatchUUID) getID(ctx context.Context) (uuid.UUID, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.size == 0 || (u.ttl != nil && u.lastTime.Add(*u.ttl).Before(time.Now())) {
		err := u.generateIDs(ctx)
		if err != nil {
			return [16]byte{}, err
		}
	}
	id := u.batchUUIDS[u.size-1]
	u.size--
	return id, nil
}

func (u *BatchUUID) String(ctx context.Context) (string, error) {
	id, err := u.getID(ctx)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (u *BatchUUID) Bytes(ctx context.Context) ([]byte, error) {
	id, err := u.getID(ctx)
	if err != nil {
		return nil, err
	}
	return []byte(id.String()), nil
}
