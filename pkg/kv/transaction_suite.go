package kv

import (
	"context"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"strconv"
	"sync"
)

type TransactionSuite struct {
	StoreSuite
	Store TransactionStore
}

func (s *TransactionSuite) SetupTest() {
	s.StoreSuite.Store = s.Store
}

func (s *TransactionSuite) TestTransaction() {
	s.NotNil(s.Store)

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	for i := 1; i < 1024; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := []byte(strconv.Itoa(i))
			val := randutil.RandByte(uint64(i))
			txn, err := s.Store.NewTransaction()
			s.Nil(err)

			v, err := txn.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)

			err = txn.Put(ctx, key, val)
			s.Nil(err)

			v, err = txn.Get(ctx, key)
			s.Nil(err)
			s.Equal(val, v)

			err = txn.Delete(ctx, key)
			s.Nil(err)

			v, err = txn.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)

			err = txn.Commit(ctx)
			s.Nil(v)
		}()
	}
	wg.Wait()
}

func (s *TransactionSuite) TestTransactionConflict() {
	s.NotNil(s.Store)

	ctx := context.Background()
	key := randutil.RandByte(64)
	val := randutil.RandByte(128)

	err := s.Store.Put(ctx, key, val)
	s.Nil(err)

	txn1, err := s.Store.NewTransaction()
	txn2, err := s.Store.NewTransaction()

	v, err := txn1.Get(ctx, key)
	s.Equal(val, v)
	s.Nil(err)

	v, err = txn2.Get(ctx, key)
	s.Equal(val, v)
	s.Nil(err)

	err = txn1.Put(ctx, key, randutil.RandByte(256))
	s.Nil(err)

	err = txn2.Put(ctx, key, randutil.RandByte(512))
	if err != nil {
		// pessimistic transaction
		err = txn1.Commit(ctx)
		s.Nil(err)

		err = txn2.Put(ctx, key, randutil.RandByte(512))
		if err != nil {
			err := txn2.Rollback()
			s.Nil(err)
		} else {
			err := txn2.Commit(ctx)
			if err != nil {
				err := txn2.Rollback()
				s.Nil(err)
			} else {
				s.T().Log("transaction not satisfied ACID")
			}
		}
	} else {
		// optimistic transaction
		err = txn1.Commit(ctx)
		s.Nil(err)

		err = txn2.Commit(ctx)
		s.NotNil(err)

		err = txn2.Rollback()
		s.Nil(err)
	}
}

func (s *TransactionSuite) TestTransactionConsistency() {
	s.NotNil(s.Store)

	ctx := context.Background()
	key := randutil.RandByte(64)

	txn, err := s.Store.NewTransaction()

	v, err := txn.Get(ctx, key)
	s.Nil(v)
	s.Equal(ErrKeyNotFound, err)

	err = s.Store.Put(ctx, key, randutil.RandByte(128))
	if err != nil {
		err := txn.Put(ctx, key, randutil.RandByte(256))
		s.Nil(err)

		err = txn.Commit(ctx)
		s.Nil(err)
	} else {
		err := txn.Put(ctx, key, randutil.RandByte(256))
		if err != nil {
			// optimistic transaction
			err = txn.Rollback()
			s.Nil(err)
		} else {
			err = txn.Commit(ctx)
			if err != nil {
				err = txn.Rollback()
				s.Nil(err)
			} else {
				s.T().Log("transaction not satisfied ACID")
			}
		}
	}
}
