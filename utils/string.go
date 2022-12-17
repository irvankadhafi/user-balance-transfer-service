package utils

import "strconv"

// StringToInt64 :nodoc:
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return i
}
