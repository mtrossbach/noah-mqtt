package misc

import "strconv"

func S2i(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		return -1
	} else {
		return i
	}
}

func ParseFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	} else {
		return 0
	}
}
