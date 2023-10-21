package entry

import (
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"os"
	"time"
)

type Entry struct {
	fullpath.FullPath             // file full path
	bucket.Bucket                 // own bucket
	Mtime             time.Time   // time of last modification
	Ctime             time.Time   // time of creation
	Mode              os.FileMode // file mode
	Mime              string      // MIME type
	Md5               []byte      // MD5
	FileSize          uint64      // file size
	Chunks                        // chunks
}

func (ent *Entry) Key() string {
	return string(ent.Bucket) + string(ent.FullPath) + entryKv
}

func EntryKey(bkt bucket.Bucket, fp fullpath.FullPath) string {
	return string(bkt) + string(fp) + entryKv
}

func (ent *Entry) IsDirectory() bool {
	return ent.Mode.IsDir()
}

func (ent *Entry) IsFile() bool {
	return ent.Mode.IsRegular()
}
