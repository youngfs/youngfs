package directory

import (
	"icesos/full_path"
	"icesos/iam"
	"os"
	"time"
)

type Inode struct {
	full_path.FullPath             //full full_path
	iam.Set                        //own set_iam
	Time               time.Time   //time of creation
	Mode               os.FileMode //file mode
}

func (inode *Inode) IsDirectory() bool {
	return inode.Mode&os.ModeDir > 0
}

func (inode *Inode) Key() string {
	return string(inode.Set) + string(inode.FullPath) + inodeKv
}
