package snippet

// 复制一个slice
func copySlice(s []int) (copyOfs []int) {
	copyOfs = append([]int(nil), s...)
}
