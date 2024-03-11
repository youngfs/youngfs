package volume

import (
	"go.uber.org/multierr"
	"os"
	"sync"
)

type FilePool struct {
	pool  sync.Pool
	files []*os.File
	mutex sync.Mutex
	path  string
}

func NewFilePool(filePath string) *FilePool {
	return &FilePool{
		path: filePath,
	}
}

func (fp *FilePool) openFile() *os.File {
	file, err := os.Open(fp.path)
	if err != nil {
		return nil
	}

	fp.mutex.Lock()
	fp.files = append(fp.files, file) // 加入追踪列表
	fp.mutex.Unlock()

	return file
}

func (fp *FilePool) Get() *os.File {
	file, ok := fp.pool.Get().(*os.File)
	if !ok || file == nil {
		file = fp.openFile() // 需要时创建新文件
	}
	return file
}

func (fp *FilePool) Put(file *os.File) {
	if file != nil {
		fp.pool.Put(file)
	}
}

func (fp *FilePool) Close() error {
	fp.mutex.Lock()
	defer fp.mutex.Unlock()

	var merr error
	for _, file := range fp.files {
		if err := file.Close(); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	fp.files = []*os.File{}
	return merr
}
