package util

type BytesSlice [][]byte

func (b BytesSlice) Len() int {
	return len(b)
}

func (b BytesSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BytesSlice) Less(i, j int) bool {
	lenI, lenJ := len(b[i]), len(b[j])
	mn := Min(lenI, lenJ)
	for k := 0; k < mn; k++ {
		if b[i][k] < b[j][k] {
			return true
		} else if b[i][k] > b[j][k] {
			return false
		}
	}
	return lenI < lenJ
}
