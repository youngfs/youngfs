package needle

import (
	"github.com/youngfs/youngfs/pkg/util"
	"strconv"
)

type Id uint64

func (id *Id) String() string {
	return strconv.FormatUint(uint64(*id), 10)
}

func (id *Id) Bytes() []byte {
	ret := make([]byte, idLen)
	util.Uint64toBytes(ret, uint64(*id))
	return ret
}
