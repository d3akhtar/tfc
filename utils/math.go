package utils

import "cmp"

func Clamp[T cmp.Ordered](num, min, max T) T {
	if num <= min {
		return min
	} else if num >= max {
		return max
	}

	return num
}
