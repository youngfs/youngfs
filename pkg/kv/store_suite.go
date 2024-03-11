package kv

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

type StoreSuite struct {
	suite.Suite
	Store Store
}

func (s *StoreSuite) TestKV() {
	s.NotNil(s.Store)

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	for i := 1; i < 1024; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := []byte(strconv.Itoa(i))
			val := randutil.RandByte(uint64(i))
			v, err := s.Store.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)

			err = s.Store.Put(ctx, key, val)
			s.Nil(err)

			v, err = s.Store.Get(ctx, key)
			s.Nil(err)
			s.Equal(val, v)

			err = s.Store.Delete(ctx, key)
			s.Nil(err)

			v, err = s.Store.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)
		}()
	}
	wg.Wait()
}

func (s *StoreSuite) TestIterator() {
	s.NotNil(s.Store)

	ctx := context.Background()
	set := make(map[string]bool)
	prefix := randutil.RandByte(64)
	type pair struct {
		key []byte
		val []byte
	}
	pairs := make([]*pair, 0)
	prefixPairs := make([]*pair, 0)
	for i := 0; i < 1024; i++ {
		var genKey func() []byte
		// i = 0: prefixFlag = true
		// i = 1: prefixFlag = false
		prefixFlag := i == 0 || (i != 1 && rand.Intn(2) == 0)
		if prefixFlag {
			genKey = func() []byte { return append(prefix, randutil.RandByte(32)...) }
		} else {
			genKey = func() []byte { return randutil.RandByte(32) }
		}
		key, val := genKey(), randutil.RandByte(128)
		for set[string(key)] {
			key = genKey()
		}
		set[string(key)] = true

		v, err := s.Store.Get(ctx, key)
		s.Nil(v)
		s.Equal(ErrKeyNotFound, err)

		err = s.Store.Put(ctx, key, val)
		s.Nil(err)

		v, err = s.Store.Get(ctx, key)
		s.Equal(val, v)
		s.Nil(err)

		pair := &pair{
			key: key, val: val,
		}
		pairs = append(pairs, pair)
		if prefixFlag {
			prefixPairs = append(prefixPairs, pair)
		}
	}

	sortPair := func(pairs []*pair) {
		sort.Slice(pairs, func(i, j int) bool {
			return bytes.Compare(pairs[i].key, pairs[j].key) <= 0
		})
	}
	sortPair(pairs)
	sortPair(prefixPairs)

	wg := &sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			prefixFlag := rand.Intn(2) == 0
			basePairs := pairs
			var itOpts []IteratorOption
			if prefixFlag {
				basePairs = prefixPairs
				itOpts = append(itOpts, WithPrefix(prefix))
			}
			p := rand.Intn(len(basePairs))
			it, err := s.Store.NewIterator(itOpts...)
			defer it.Close()
			s.Nil(err)

			for it.Seek(basePairs[p].key); it.Valid(); it.Next() {
				s.Equal(basePairs[p].key, it.Key())
				s.Equal(basePairs[p].val, it.Value())
				p++
			}

			s.Equal(len(basePairs), p)
		}()
	}
	wg.Wait()
}

type TTLStoreSuite struct {
	suite.Suite
	Store Store
	TTL   time.Duration
}

func (s *TTLStoreSuite) TestKV() {
	s.NotNil(s.Store)

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	for i := 1; i < 1024; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := []byte(strconv.Itoa(i))
			val := randutil.RandByte(uint64(i))
			v, err := s.Store.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)

			err = s.Store.Put(ctx, key, val)
			s.Nil(err)

			time.Sleep(s.TTL / 2)

			v, err = s.Store.Get(ctx, key)
			s.Nil(err)
			s.Equal(val, v)

			if (i & 1) == 0 {
				err = s.Store.Delete(ctx, key)
				s.Nil(err)

				v, err = s.Store.Get(ctx, key)
				s.Equal(ErrKeyNotFound, err)
				s.Nil(v)
			}
			time.Sleep(s.TTL)

			v, err = s.Store.Get(ctx, key)
			s.Equal(ErrKeyNotFound, err)
			s.Nil(v)
		}()
	}
	wg.Wait()
}
