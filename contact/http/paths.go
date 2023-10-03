package http

import (
	"cmp"
	"fmt"
	"slices"
)

type ResourcePaths struct {
	Root, Form, List string
}

type validResourcePaths ResourcePaths

// paths should be distinct or this will panic
func (my ResourcePaths) Validated() (v validResourcePaths, err error) {
	if hasDuplicates([]string{my.Root, my.Form, my.List}) {
		return v, fmt.Errorf("path elements must be unique. Got %+v", my)
	}
	return validResourcePaths(my), nil
}

func hasDuplicates[T cmp.Ordered](s []T) bool {
	initialLen := len(s)
	slices.Sort(s)
	s = slices.Compact(s)
	return len(s) != initialLen
}
