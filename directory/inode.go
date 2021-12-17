package directory

import (
	"icesos/directory/directory_pb"
	"icesos/full_path"
	"icesos/iam"
	"os"
	"time"
)

type Inode struct {
	full_path.FullPath             //full full_path
	set_iam.Set                    //own set_iam
	Time               time.Time   //time of creation
	Mode               os.FileMode //file mode
}

func (inode *Inode) IsDirectory() bool {
	return inode.Mode&os.ModeDir > 0
}

func (inode *Inode) TimeUnix() int64 {
	return inode.Time.Unix()
}

func (inode *Inode) Key() string {
	return string(inode.Set) + "_" + string(inode.FullPath) + "_inode"
}

func (inode *Inode) toPb() *directory_pb.Inode {
	if inode == nil {
		return nil
	}
	return &directory_pb.Inode{
		FullPath: string(inode.FullPath),
		Set:      string(inode.Set),
		Time:     inode.Time.Unix(),
		Mode:     uint32(inode.Mode),
	}
}

func inodePbToInstance(pb *directory_pb.Inode) *Inode {
	if pb == nil {
		return nil
	}
	return &Inode{
		FullPath: full_path.FullPath(pb.FullPath),
		Set:      set_iam.Set(pb.Set),
		Time:     time.Unix(pb.Time, 0),
		Mode:     os.FileMode(pb.Mode),
	}
}
