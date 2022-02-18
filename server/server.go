package server

import (
	"icesos/filer"
	"icesos/storage_engine"
)

type Server struct {
	Filer         *filer.Filer
	StorageEngine *storage_engine.StorageEngine
}
