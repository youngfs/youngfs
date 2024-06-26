package entry

import (
	"github.com/golang/protobuf/proto"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/entry/entry_pb"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"os"
	"time"
)

func (ent *Entry) toPb() *entry_pb.Entry {
	if ent == nil {
		return nil
	}

	chunks := make([]*entry_pb.Chunk, len(ent.Chunks))
	for i, u := range ent.Chunks {
		chunks[i] = u.toPb()
	}

	mtime, ctime := int64(0), int64(0)
	if !ent.IsDirectory() {
		mtime, ctime = ent.Mtime.UnixNano(), ent.Ctime.UnixNano()
	}

	return &entry_pb.Entry{
		FullPath: string(ent.FullPath),
		Bucket:   string(ent.Bucket),
		Mtime:    mtime,
		Ctime:    ctime,
		Mode:     uint32(ent.Mode),
		Mine:     ent.Mime,
		Md5:      ent.Md5,
		FileSize: ent.FileSize,
		Chunks:   chunks,
	}
}

func (c *Chunk) toPb() *entry_pb.Chunk {
	if c == nil {
		return nil
	}

	frags := make([]*entry_pb.Frag, len(c.Frags))
	for i, u := range c.Frags {
		frags[i] = u.toPb()
	}

	return &entry_pb.Chunk{
		Offset:        c.Offset,
		Size:          c.Size,
		Md5:           c.Md5,
		IsReplication: c.IsReplication,
		Frags:         frags,
	}
}

func (f *Frag) toPb() *entry_pb.Frag {
	if f == nil {
		return nil
	}

	return &entry_pb.Frag{
		Size:        f.Size,
		Id:          f.Id,
		Md5:         f.Md5,
		IsDataShard: f.IsDataShard,
		Fid:         f.Fid,
	}
}

func entryPbToInstance(pb *entry_pb.Entry) *Entry {
	if pb == nil {
		return nil
	}

	chunks := make([]*Chunk, len(pb.Chunks))
	for i, u := range pb.Chunks {
		if u == nil {
			continue
		}
		chunks[i] = chunkPbToInstance(u)
	}

	if pb.Chunks == nil {
		chunks = nil
	}

	var mtime, ctime time.Time
	mode := os.FileMode(pb.Mode)
	if !mode.IsDir() {
		mtime = time.Unix(pb.Mtime/int64(time.Second), pb.Mtime%int64(time.Second))
		ctime = time.Unix(pb.Ctime/int64(time.Second), pb.Ctime%int64(time.Second))
	}

	return &Entry{
		FullPath: fullpath.FullPath(pb.FullPath),
		Bucket:   bucket.Bucket(pb.Bucket),
		Mtime:    mtime,
		Ctime:    ctime,
		Mode:     mode,
		Mime:     pb.Mine,
		Md5:      pb.Md5,
		FileSize: pb.FileSize,
		Chunks:   chunks,
	}
}

func chunkPbToInstance(pb *entry_pb.Chunk) *Chunk {
	if pb == nil {
		return nil
	}

	frags := make([]*Frag, len(pb.Frags))
	for i, u := range pb.Frags {
		if u == nil {
			continue
		}
		frags[i] = frgaPbToInstance(u)
	}

	if pb.Frags == nil {
		frags = nil
	}

	return &Chunk{
		Offset:        pb.Offset,
		Size:          pb.Size,
		Md5:           pb.Md5,
		IsReplication: pb.IsReplication,
		Frags:         frags,
	}
}

func frgaPbToInstance(pb *entry_pb.Frag) *Frag {
	if pb == nil {
		return nil
	}

	return &Frag{
		Size:        pb.Size,
		Id:          pb.Id,
		Md5:         pb.Md5,
		IsDataShard: pb.IsDataShard,
		Fid:         pb.Fid,
	}
}

func (ent *Entry) EncodeProto() ([]byte, error) {
	message := ent.toPb()
	b, err := proto.Marshal(message)
	if err != nil {
		err = errors.ErrProto.WarpErr(err)
	}
	return b, err
}

func DecodeEntryProto(b []byte) (*Entry, error) {
	message := &entry_pb.Entry{}
	if err := proto.Unmarshal(b, message); err != nil {
		return nil, errors.ErrProto.WarpErr(err)
	}
	return entryPbToInstance(message), nil
}
