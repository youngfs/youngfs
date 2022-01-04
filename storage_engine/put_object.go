package storage_engine

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}
