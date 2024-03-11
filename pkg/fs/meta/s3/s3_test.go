package s3

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/kv/badger"
	"github.com/youngfs/youngfs/pkg/util/randutil"
	"math/rand"
	"os"
	"path"
	"sort"
	"testing"
	"time"
)

type s3Suite struct {
	suite.Suite
	s3 *S3
	kv []kv.TransactionStore
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(s3Suite))
}

func (s *s3Suite) SetupTest() {
	dir := s.T().TempDir()
	s3Kv, err := badger.New(path.Join(dir, "s3kv"))
	s.Nil(err)
	s.kv = append(s.kv, s3Kv)
	continueKv, err := badger.New(path.Join(dir, "s3continuekv"), badger.WithTTL(24*time.Hour))
	s.Nil(err)
	s.kv = append(s.kv, continueKv)
	s.s3 = New(s3Kv, continueKv)
}

func (s *s3Suite) TearDownTest() {
	for _, store := range s.kv {
		err := store.Close()
		s.Nil(err)
	}
}

func (s *s3Suite) TestSimpleObject() {
	s.NotNil(s.s3)

	bkt := bucket.Bucket("test")
	now := time.Now()
	ctx := context.Background()
	nowTime := time.Unix(now.UnixNano()/int64(time.Second), now.UnixNano()%int64(time.Second))
	ents := []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/cc"},
		{Bucket: bkt, FullPath: "/aa/bc/dd/ee"},
		{Bucket: bkt, FullPath: "/aa/cc/dd"},
		{Bucket: bkt, FullPath: "/aa/cc/ee/ff"},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/ff"},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range ents {
		ent.Mtime = nowTime
		ent.Ctime = nowTime
		ent.FileSize = 1024
		ent.Mode = os.ModePerm
	}

	err := s.s3.S3PutObjects(ctx, ents)
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.Nil(err)
		s.Equal(ents[i], ent)
	}

	ret, err := s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		ContinueToken: randutil.RandString(16),
	})
	s.ErrorIs(err, errors.ErrInvalidContinueToken)
	s.Nil(ret)

	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix: entry.EntryKey(bkt, "/"),
	})
	s.Nil(err)
	s.Equal(ents, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(ents), ret.KeyCount)

	expectedEnts := []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/"),
		Delimiter: "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/cc"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/aa/bb/cc"),
		Delimiter: "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/bc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/aa/"),
		Delimiter: "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/cc"},
		{Bucket: bkt, FullPath: "/aa/bc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/ff"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/aa/"),
		Delimiter: "c/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bc/dd/ee"},
		{Bucket: bkt, FullPath: "/aa/cc/dd"},
		{Bucket: bkt, FullPath: "/aa/cc/ee/ff"},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/ff"},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:     entry.EntryKey(bkt, "/"),
		StartAfter: entry.EntryKey(bkt, "/aa/bc/"),
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:     entry.EntryKey(bkt, "/aa/"),
		StartAfter: entry.EntryKey(bkt, "/aa/cc/"),
		Delimiter:  "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	err = s.s3.S3DeleteObjects(ctx, []string{
		entry.EntryKey(bkt, "/aa/bb/cc"),
		entry.EntryKey(bkt, "/aa/bc/dd/ee"),
		entry.EntryKey(bkt, "/aa/cc/dd"),
		entry.EntryKey(bkt, "/aa/cc/ee/ff"),
		entry.EntryKey(bkt, "/aa/dd"),
		entry.EntryKey(bkt, "/aa/ee/ff"),
		entry.EntryKey(bkt, "/bb"),
	})
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.ErrorIs(err, errors.ErrObjectNotFound)
		s.Nil(ent)
	}
}

func (s *s3Suite) TestDirObject() {
	s.NotNil(s.s3)

	bkt := bucket.Bucket("test")
	now := time.Now()
	ctx := context.Background()
	nowTime := time.Unix(now.UnixNano()/int64(time.Second), now.UnixNano()%int64(time.Second))
	ents := []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/cc"},
		{Bucket: bkt, FullPath: "/aa/bc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/bc/dd/ee"},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/dd"},
		{Bucket: bkt, FullPath: "/aa/cc/ee/ff"},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range ents {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}

	err := s.s3.S3PutObjects(ctx, ents)
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.Nil(err)
		s.Equal(ents[i], ent)
	}

	ret, err := s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:        entry.EntryKey(bkt, "/"),
		ContinueToken: randutil.RandString(16),
	})
	s.ErrorIs(err, errors.ErrInvalidContinueToken)
	s.Nil(ret)

	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{})
	s.Nil(err)
	s.Equal(ents, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(ents), ret.KeyCount)

	expectedEnts := []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/"),
		Delimiter: "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/bc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/aa/"),
		Delimiter: "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bb/cc"},
		{Bucket: bkt, FullPath: "/aa/bc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:    entry.EntryKey(bkt, "/aa/"),
		Delimiter: "c/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/bc/dd/ee"},
		{Bucket: bkt, FullPath: "/aa/cc/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/aa/cc/dd"},
		{Bucket: bkt, FullPath: "/aa/cc/ee/ff"},
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
		{Bucket: bkt, FullPath: "/bb"},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:     entry.EntryKey(bkt, "/"),
		StartAfter: entry.EntryKey(bkt, "/aa/bc/"),
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	expectedEnts = []*entry.Entry{
		{Bucket: bkt, FullPath: "/aa/dd"},
		{Bucket: bkt, FullPath: "/aa/ee/", Mode: os.ModeDir},
	}
	for _, ent := range expectedEnts {
		if ent.Mode != os.ModeDir {
			ent.Mode = os.ModePerm
			ent.Mtime = nowTime
			ent.Ctime = nowTime
			ent.FileSize = 1024
		}
	}
	ret, err = s.s3.S3ListObjects(ctx, &ListObjectsOptions{
		Prefix:     entry.EntryKey(bkt, "/aa/"),
		StartAfter: entry.EntryKey(bkt, "/aa/cc/"),
		Delimiter:  "/",
	})
	s.Nil(err)
	s.Equal(expectedEnts, ret.Entries)
	s.Equal(false, ret.IsTruncated)
	s.Equal("", ret.NextContinuationToken)
	s.Equal(len(expectedEnts), ret.KeyCount)

	err = s.s3.S3DeleteObjects(ctx, []string{
		entry.EntryKey(bkt, "/aa/bb/cc"),
		entry.EntryKey(bkt, "/aa/bc/"),
		entry.EntryKey(bkt, "/aa/bc/dd/ee"),
		entry.EntryKey(bkt, "/aa/cc/"),
		entry.EntryKey(bkt, "/aa/cc/dd"),
		entry.EntryKey(bkt, "/aa/cc/ee/ff"),
		entry.EntryKey(bkt, "/aa/dd"),
		entry.EntryKey(bkt, "/aa/ee/"),
		entry.EntryKey(bkt, "/bb"),
	})
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.ErrorIs(err, errors.ErrObjectNotFound)
		s.Nil(ent)
	}
}

func (s *s3Suite) TestLargeObjects() {
	s.NotNil(s.s3)

	bkt := bucket.Bucket("test")
	now := time.Now()
	ctx := context.Background()
	nowTime := time.Unix(now.UnixNano()/int64(time.Second), now.UnixNano()%int64(time.Second))
	ents := make([]*entry.Entry, 0)
	size := 128 * 1024
	for i := 0; i < size; i++ {
		ents = append(ents, &entry.Entry{
			Bucket:   bkt,
			FullPath: fullpath.FullPath(fmt.Sprintf("/%d.txt", i)),
			Mtime:    nowTime,
			Ctime:    nowTime,
			FileSize: 1024,
			Mode:     os.ModePerm,
		})
	}
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].FullPath < ents[j].FullPath
	})

	err := s.s3.S3PutObjects(ctx, ents)
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.Nil(err)
		s.Equal(ents[i], ent)
	}

	for o := 0; o < 4; o++ {
		for from, cnt, continueToken := 0, 0, ""; from < size; from += cnt {
			opt := &ListObjectsOptions{
				Prefix:        entry.EntryKey(bkt, "/"),
				MaxKeys:       rand.Intn(64),
				ContinueToken: continueToken,
			}
			cnt = opt.MaxKeys
			if cnt == 0 {
				cnt = 1000
			}
			cnt = min(cnt, size-from)
			ret, err := s.s3.S3ListObjects(ctx, opt)
			s.Nil(err)
			s.Equal(cnt, ret.KeyCount)
			s.Equal(ents[from:from+cnt], ret.Entries)
			s.Equal(from+cnt != size, ret.IsTruncated)
			if ret.IsTruncated {
				s.NotEqual("", ret.NextContinuationToken)
			}
			continueToken = ret.NextContinuationToken
		}
	}

	for o := 0; o < 4; o++ {
		for from, cnt, startAfter := 0, 0, ""; from < size; from += cnt {
			opt := &ListObjectsOptions{
				Prefix:     entry.EntryKey(bkt, "/"),
				MaxKeys:    rand.Intn(10),
				StartAfter: startAfter,
			}
			cnt = opt.MaxKeys
			if cnt == 0 {
				cnt = 1000
			}
			cnt = min(cnt, size-from)
			ret, err := s.s3.S3ListObjects(ctx, opt)
			s.Nil(err)
			s.Equal(cnt, ret.KeyCount)
			s.Equal(ents[from:from+cnt], ret.Entries)
			s.Equal(from+cnt != size, ret.IsTruncated)
			if ret.IsTruncated {
				s.NotEqual("", ret.NextContinuationToken)
			}
			startAfter = entry.EntryKey(bkt, ret.Entries[cnt-1].FullPath)
		}
	}

	objs := make([]string, 0)
	for _, ent := range ents {
		objs = append(objs, entry.EntryKey(bkt, ent.FullPath))
	}
	err = s.s3.S3DeleteObjects(ctx, objs)
	s.Nil(err)

	for i, _ := range ents {
		ent, err := s.s3.S3GetObject(ctx, entry.EntryKey(bkt, ents[i].FullPath))
		s.ErrorIs(err, errors.ErrObjectNotFound)
		s.Nil(ent)
	}
}
