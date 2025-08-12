package common

import "strconv"

func AtoiWithFallback(s string, fallback int) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}

	return res
}
