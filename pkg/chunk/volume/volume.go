package volume

import (
	"fmt"
	"github.com/youngfs/youngfs/pkg/chunk/volume/needle"
	"github.com/youngfs/youngfs/pkg/errors"
	"go.uber.org/multierr"
	"io"
	"os"
	"path"
)

type Volume struct {
	dir         string
	reader      *FilePool
	writer      *os.File
	id          uint64
	count       uint64
	fileSize    uint64
	needleStore needle.Store
	magic       []byte
}

func NewVolume(dir string, id uint64, creator needle.StoreCreator) (*Volume, error) {
	writer, err := os.OpenFile(path.Join(dir, fmt.Sprintf("%d.data")), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	reader := NewFilePool(path.Join(dir, fmt.Sprintf("%d.data")))
	needleStore, err := creator(path.Join(dir, fmt.Sprintf("%d.idx")))
	if err != nil {
		return nil, err
	}
	stat, err := writer.Stat()
	if err != nil {
		return nil, err
	}
	return &Volume{
		reader:      reader,
		writer:      writer,
		id:          id,
		count:       needleStore.Size(),
		fileSize:    uint64(stat.Size()),
		needleStore: needleStore,
	}, nil
}

func (v *Volume) Write(reader io.Reader) (string, error) {
	n, err := io.Copy(v.writer, reader)
	if err != nil {
		return "", errors.ErrVolumeWrite.WarpErr(err)
	}
	v.count += 1
	v.fileSize += uint64(n)
	nl := needle.Id(v.count)
	err = v.needleStore.Put(&needle.Needle{Id: nl, Offset: v.fileSize - uint64(n), Size: uint64(n)})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d-%s", v.id, nl.String()), nil
}

func (v *Volume) Close() error {
	var merr error
	if err := v.reader.Close(); err != nil {
		merr = multierr.Append(merr, err)
	}
	if err := v.writer.Close(); err != nil {
		merr = multierr.Append(merr, err)
	}
	if err := v.needleStore.Close(); err != nil {
		merr = multierr.Append(merr, err)
	}
	return merr
}

func (v *Volume) Size() uint64 {
	return v.fileSize
}

func (v *Volume) Magic() []byte {
	return v.magic
}

func (v *Volume) WriteMagic(magic []byte) error {
	_, err := v.writer.Write(magic)
	if err != nil {
		return errors.ErrVolumeWrite.WarpErr(err)
	}
	v.fileSize += uint64(len(magic))
	return nil
}

func (v *Volume) Read(id needle.Id, writer io.Writer) error {
	nd, err := v.needleStore.Get(id)
	if err != nil {
		return err
	}
	f := v.reader.Get()
	defer v.reader.Put(f)
	_, err = f.Seek(int64(nd.Offset), io.SeekStart)
	if err != nil {
		return errors.ErrVolumeRead.WarpErr(err)
	}
	_, err = io.CopyN(writer, f, int64(nd.Size))
	if err != nil {
		return errors.ErrVolumeRead.WarpErr(err)
	}
	return nil
}

func (v *Volume) Delete(id needle.Id) error {
	err := v.needleStore.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
