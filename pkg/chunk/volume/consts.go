package volume

import "time"

const (
	maxVolumeSize     = 1024 * 1024 * 1024 * 16 // 16GiB
	maxWritableVolume = 8
	magicLen          = 64
	heartbeatTick     = 5 * time.Second
)

const (
	numFileName = "num.lock"
)
