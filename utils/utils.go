package utils

import "strconv"

func ConvertStringToInt(s string) int {
	value , err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return value
}
