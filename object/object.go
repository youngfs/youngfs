package object

import (
	"crypto/md5"
	"icesos/full_path"
	"icesos/iam"
	"icesos/kv"
	"os"
	"time"
)

type Object struct {
	full_path.FullPath                // file full full_path
	iam.Set                           // own set_iam
	Time               time.Time      // time of creation
	Mode               os.FileMode    // file mode
	Mime               string         // MIME type
	Md5                [md5.Size]byte // MD5
	FileSize           uint64         // file size
	VolumeId           uint64         // volume id
	Fid                string         // fid
}

func (ob *Object) key() string {
	return string(ob.Set) + "_" + string(ob.FullPath) + objectKv
}

func objectKey(fp full_path.FullPath, set iam.Set) string {
	return string(set) + "_" + string(fp) + objectKv
}

func PutObject(ob *Object) error {

	return nil
}

func GetObject(fp full_path.FullPath, set iam.Set) (*Object, error) {
	key := objectKey(fp, set)

	b, err := kv.Client.KvGet(key)
	if err != nil {
		return nil, err
	}

	return decodeObjectProto(b)
}
