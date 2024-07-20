package cronparser

import (
	"strconv"
	"strings"
)

func buildIntList(min, max, interval int) []int {
	var intList []int
	for i := min; i <= max; i += interval {
		intList = append(intList, i)
	}

	return intList
}

func intsJoin(ints []int, sep string) string {
	strInts := make([]string, len(ints))
	for i, v := range ints {
		strInts[i] = strconv.Itoa(v)
	}

	return strings.Join(strInts, sep)
}
