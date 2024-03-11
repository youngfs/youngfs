package s3

import (
	"context"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
)

func (s *S3) InsertObject(ctx context.Context, ent *entry.Entry) error {
	if err := s.DeleteObjects(ctx, ent.Bucket, ent.FullPath); err != nil {
		return err
	}
	return s.S3PutObject(ctx, ent)
}

func (s *S3) GetObject(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) (*entry.Entry, error) {
	return s.S3GetObject(ctx, entry.EntryKey(bkt, fp))
}

func (s *S3) ListObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath, recursive bool) ([]*entry.Entry, error) {
	ret := make([]*entry.Entry, 0)
	delimiter := "/"
	if recursive {
		delimiter = ""
	}
	for nextContinueToken, isTruncated := "", true; isTruncated; {
		lret, err := s.S3ListObjects(ctx, &ListObjectsOptions{
			Prefix:        entry.EntryKey(bkt, fp),
			ContinueToken: nextContinueToken,
			Delimiter:     delimiter,
		})
		if err != nil {
			return nil, err
		}
		ret = append(ret, lret.Entries...)
		nextContinueToken = lret.NextContinuationToken
		isTruncated = lret.IsTruncated
	}
	return ret, nil
}

func (s *S3) DeleteObjects(ctx context.Context, bkt bucket.Bucket, fp fullpath.FullPath) error {
	ents, err := s.ListObjects(ctx, bkt, fp, true)
	if err != nil {
		return err
	}
	if len(ents) > 0 {
		objs := make([]string, len(ents))
		for i, e := range ents {
			objs[i] = e.Key()
		}
		if err := s.S3DeleteObjects(ctx, objs); err != nil {
			return err
		}
	}
	return nil
}
