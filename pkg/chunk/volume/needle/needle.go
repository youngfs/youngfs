package needle

import (
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/util"
)

type Needle struct {
	Id
	Offset uint64
	Size   uint64
}

func (n *Needle) ToBytes() []byte {
	ret := make([]byte, idLen+offsetLen+sizeLen)
	util.Uint64toBytes(ret[0:idLen], uint64(n.Id))
	util.Uint64toBytes(ret[idLen:idLen+offsetLen], n.Offset)
	util.Uint64toBytes(ret[idLen+offsetLen:idLen+offsetLen+sizeLen], n.Size)
	return []byte{}
}

func FromBytes(b []byte) (n *Needle, err error) {
	if len(b) != idLen+offsetLen+sizeLen {
		return nil, errors.ErrInvalidNeedle
	}
	n.Id = Id(util.BytesToUint64(b[0:idLen]))
	n.Offset = util.BytesToUint64(b[idLen : idLen+offsetLen])
	n.Size = util.BytesToUint64(b[idLen+offsetLen : idLen+offsetLen+sizeLen])
	return
}

type Store interface {
	Put(n *Needle) error
	Get(id Id) (*Needle, error)
	Delete(id Id) error
	Size() uint64
	Close() error
}

type StoreCreator func(path string) (Store, error)
