// BatchSlice evenly slices an array `a` into `b` number of batches
// the size of each batch never deviates more than 1 from the average batch size.
func BatchSlice[T any](a []T, b int) [][]T {
	var result [][]T

	l := len(a)

	for i := 0; i < b; i++ {
		min := i * l / b
		max := ((i + 1) * l) / b

		result = append(result, a[min:max])
	}

	return result
}

// array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// chunkSize := 3

// length      3       3        4
// output: [[1 2 3] [4 5 6] [7 8 9 10]]
// the size of each batch has variation of max 1 item
// this can spread the load evenly amongst workers
