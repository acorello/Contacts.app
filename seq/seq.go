package seq

import (
	"cmp"
	"slices"
)

func Map[T any, U any](ts []T, f func(T) U) (res []U) {
	if ts == nil {
		return nil
	}
	for _, v := range ts {
		res = append(res, f(v))
	}
	return
}

func HasDuplicates[T cmp.Ordered](s ...T) bool {
	initialLen := len(s)
	slices.Sort(s)
	s = slices.Compact(s)
	return len(s) != initialLen
}
