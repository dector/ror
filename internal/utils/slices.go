package utils

func SubSlice(slice []string, from int) []string {
	if from <= 0 {
		return slice
	}

	if from >= len(slice) {
		return slice
	}

	return slice[from:]
}
