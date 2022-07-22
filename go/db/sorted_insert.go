package db

import "sort"

// https://stackoverflow.com/questions/42746972/golang-sortedInsert-to-a-sorted-slice
func sortedInsert(ss []string, s string) []string {
	i := sort.SearchStrings(ss, s)
	ss = append(ss, "")
	copy(ss[i+1:], ss[i:])
	ss[i] = s
	return ss
}
