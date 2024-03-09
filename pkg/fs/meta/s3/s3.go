package s3

import (
	"bytes"
	"context"
	"errors"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/idgenerator"
	"github.com/youngfs/youngfs/pkg/idgenerator/uuid"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/util"
	"os"
)

const (
	batchSize = 1024 * 16
)

type S3 struct {
	kv         kv.TransactionStore
	continueKv kv.TransactionStore
	idGen      idgenerator.IDGenerator
}

func New(kv, continueKv kv.TransactionStore, options ...Option) *S3 {
	cfg := &config{
		idGen: uuid.NewSimpleUUID(),
	}
	for _, opt := range options {
		opt.apply(cfg)
	}
	return &S3{
		kv:         kv,
		continueKv: continueKv,
		idGen:      cfg.idGen,
	}
}

func (s *S3) S3PutObject(ctx context.Context, ent *entry.Entry) error {
	val, err := ent.EncodeProto()
	if err != nil {
		return err
	}
	return s.kv.Put(ctx, []byte(ent.Key()), val)
}

func (s *S3) S3PutObjects(ctx context.Context, entries []*entry.Entry) error {
	for p := 0; p < len(entries); p += batchSize {
		err := kv.DoTransaction(s.kv, ctx, func(txn kv.Transaction) error {
			up := min(p+batchSize, len(entries))
			for i := p; i < up; i++ {
				val, err := entries[i].EncodeProto()
				if err != nil {
					return err
				}

				err = txn.Put(ctx, []byte(entries[i].Key()), val)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *S3) S3GetObject(ctx context.Context, obj string) (*entry.Entry, error) {
	val, err := s.kv.Get(ctx, []byte(obj))
	if err != nil {
		if errors.Is(err, kv.ErrKeyNotFound) {
			return nil, ErrS3ObjectNotFound
		} else {
			return nil, err
		}
	}

	return entry.DecodeEntryProto(val)
}

func (s *S3) S3DeleteObject(ctx context.Context, obj string) error {
	return s.kv.Delete(ctx, []byte(obj))
}

func (s *S3) S3DeleteObjects(ctx context.Context, objs []string) error {
	for p := 0; p < len(objs); p += batchSize {
		err := kv.DoTransaction(s.kv, ctx, func(txn kv.Transaction) error {
			up := min(p+batchSize, len(objs))
			for i := p; i < up; i++ {
				err := txn.Delete(ctx, []byte(objs[i]))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type ListObjectsOptions struct {
	Prefix        string
	Delimiter     string
	ContinueToken string
	// StartAfter is where you should start if you want to fetch the next page of results.
	StartAfter string
	MaxKeys    int
}

type ListObjectsResult struct {
	Entries []*entry.Entry `json:"entries"`
	// NextContinuationToken is the token to use to get the next page of results.
	NextContinuationToken string `json:"nextContinuationToken"`
	// IsTruncated is whether the listing is truncated.
	IsTruncated bool `json:"isTruncated"`
	// keyCount is the number of keys in the listing.
	KeyCount int `json:"keyCount"`
}

func (s *S3) S3ListObjects(ctx context.Context, opt *ListObjectsOptions) (*ListObjectsResult, error) {
	if opt == nil {
		return nil, errors.New("list objects options is nil")
	}
	prefix := []byte(opt.Prefix)
	delimiter := []byte(opt.Delimiter)
	if opt.MaxKeys == 0 {
		opt.MaxKeys = 1000
	}
	opt.MaxKeys = min(opt.MaxKeys, 1000)
	it, err := s.kv.NewIterator(kv.WithPrefix(prefix))
	if err != nil {
		return nil, err
	}
	defer it.Close()

	startKey := prefix
	if opt.StartAfter != "" {
		startKey = []byte(opt.StartAfter)
	} else if opt.ContinueToken != "" {
		startAfter, err := s.GetContinueToken(ctx, opt.Prefix, opt.ContinueToken)
		if err != nil {
			return nil, err
		}
		startKey = []byte(startAfter)
	}
	keyCnt := 0
	entries := make([]*entry.Entry, 0)
	valid := true     // check next common prefix is valid
	it.Seek(startKey) // jump to start key
	if opt.StartAfter != "" && it.Valid() {
		key := it.Key()
		if delimIndex := bytes.Index(key[len(prefix):], delimiter); len(delimiter) != 0 && delimIndex != -1 {
			commonPrefix := key[:len(prefix)+delimIndex+len(delimiter)]
			if bytes.Equal(commonPrefix, startKey) {
				nxt := util.PrefixEnd(commonPrefix)
				if nxt == nil {
					valid = false
				} else {
					it.Seek(nxt)
				}
			}
		} else if bytes.Equal(key, startKey) {
			it.Next()
		}
	}
	if !valid {
		goto end
	}
	for it.Valid() {
		if len(delimiter) != 0 {
			key := it.Key()
			// Because the delimiter may appear in the prefix, the query is key[len(prefix):]
			if delimIndex := bytes.Index(key[len(prefix):], delimiter); delimIndex != -1 {
				commonPrefix := key[:len(prefix)+delimIndex+len(delimiter)]
				ent := entry.New(string(commonPrefix))
				ent.Mode = os.ModeDir
				entries = append(entries, ent)
				keyCnt++
				nxt := util.PrefixEnd(commonPrefix)
				if nxt == nil {
					valid = false
					break
				}
				it.Seek(nxt)
				if keyCnt == opt.MaxKeys {
					break
				}
				continue
			}
		}
		ent, err := entry.DecodeEntryProto(it.Value())
		if err != nil {
			return nil, err
		}
		entries = append(entries, ent)
		keyCnt++
		if keyCnt == opt.MaxKeys {
			break
		}
		it.Next()
	}
end:
	isTruncated := false
	continueToken := ""
	if valid && it.Valid() {
		if it.Next(); it.Valid() {
			isTruncated = true
			nextKey := it.Key()
			continueToken, err = s.NewContinueToken(ctx, opt.Prefix, string(nextKey))
			if err != nil {
				return nil, err
			}
		}
	}
	return &ListObjectsResult{
		Entries:               entries,
		IsTruncated:           isTruncated,
		NextContinuationToken: continueToken,
		KeyCount:              keyCnt,
	}, err
}
