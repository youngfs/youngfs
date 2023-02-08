package server

const (
	partSize            = 64 * 1024 * 1024
	smallObjectSize     = 4 * 1024
	replicationNum      = 2
	reedSolomonMaxShard = 6
)

var dataShardsPlan = [][]bool{
	{},
	{true},
	{true, false},
	{true, true, false},
	{true, true, true, false},
	{true, true, true, false, false},
	{true, true, true, true, false, false},
}

var dataParityShards = [][2]int{
	{1, 0},
	{1, 0},
	{1, 0},
	{2, 1},
	{3, 1},
	{3, 2},
	{4, 2},
}
