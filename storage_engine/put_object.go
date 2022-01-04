package storage_engine

type putObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}
