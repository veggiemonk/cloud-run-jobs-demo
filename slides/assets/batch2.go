array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
chunkSize := 3
var result [][]int

for i := 0; i < len(array); i += chunkSize {
	end := i + chunkSize

	if end > len(array) {
		end = len(array)
	}

	result = append(result, array[i:end])
}

fmt.Println(result)
// length       4    |    4    |  2 
// output: [[1 2 3 4] [5 6 7 8] [9 10]]
// 2 workers will do double the work of the last worker.
}
