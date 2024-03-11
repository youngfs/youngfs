package volume

import (
	"os"
	"strconv"
	"sync"
)

type NumberFile struct {
	Path string
	mux  *sync.Mutex
}

func NewNumberFile(path string) *NumberFile {
	return &NumberFile{Path: path}
}

func (nf *NumberFile) WriteNumber(num uint64) error {
	numStr := strconv.FormatUint(num, 10)
	return os.WriteFile(nf.Path, []byte(numStr), 0644)
}

func (nf *NumberFile) ReadNumber() (uint64, error) {
	data, err := os.ReadFile(nf.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	return strconv.ParseUint(string(data), 10, 64)
}

func (nf *NumberFile) Add(delta uint64) (uint64, error) {
	nf.mux.Lock()
	defer nf.mux.Unlock()
	num, err := nf.ReadNumber()
	if err != nil {
		return 0, err
	}
	num += delta
	err = nf.WriteNumber(num)
	return num, err
}

func (nf *NumberFile) WriteMax(num uint64) (uint64, error) {
	nf.mux.Lock()
	defer nf.mux.Unlock()
	old, err := nf.ReadNumber()
	if err != nil {
		return 0, err
	}
	if old < num {
		err = nf.WriteNumber(num)
	}
	return max(num, old), err
}
