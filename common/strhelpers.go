package common

import "strconv"

func ParseIntWithFallback(s string, fallback int) int {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fallback
	}

	return int(res)
}

func ParseUIntWithFallback(s string, fallback uint) uint {
	res, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fallback
	}

	return uint(res)
}
