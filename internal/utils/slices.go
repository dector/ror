package utils

func SubSlice(slice []string, from int) []string {
	if from <= 0 {
		return slice
	}

	switch len(slice) {
	case 1:
		return []string{}
	default:
		return slice[1:]
	}
}
