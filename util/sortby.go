package util

import (
	"sort"
)

type sortInterface struct {
	length int
	less   func(int, int) bool
	swap   func(int, int)
}

func (si *sortInterface) Len() int {
	return si.length
}

func (si *sortInterface) Less(i, j int) bool {
	return si.less(i, j)
}

func (si *sortInterface) Swap(i, j int) {
	si.swap(i, j)
}

func SortBy(length int, less func(int, int) bool, swap func(int, int)) {
	sort.Sort(&sortInterface{length, less, swap})
}
