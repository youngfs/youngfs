package directory

type InodeSlice []Inode

func (inodes InodeSlice) Len() int {
	return len(inodes)
}

func (inodes InodeSlice) Swap(i, j int) {
	inodes[i], inodes[j] = inodes[j], inodes[i]
}

func (inodes InodeSlice) Less(i, j int) bool {
	return inodes[i].FullPath < inodes[j].FullPath
}
