package http

import (
	"cmp"
	"fmt"
	"html/template"
	"net/url"
	"slices"

	"dev.acorello.it/go/contacts/contact"
)

type ResourcePaths struct {
	Root, Form, List string
}

type validResourcePaths ResourcePaths

// paths should be distinct or this will panic
func (my ResourcePaths) Validated() (v validResourcePaths, err error) {
	if hasDuplicates(my.Root, my.Form, my.List) {
		return v, fmt.Errorf("path elements must be unique. Got %+v", my)
	}
	return validResourcePaths(my), nil
}

func (my validResourcePaths) ContactFormURL(c contact.Contact) template.URL {
	res, err := url.Parse(my.Form)
	if err != nil {
		panic(err)
	}
	q := url.Values{}
	q.Add("Id", c.Id.String())
	res.RawQuery = q.Encode()
	return template.URL(res.String())
}

func hasDuplicates[T cmp.Ordered](s ...T) bool {
	initialLen := len(s)
	slices.Sort(s)
	s = slices.Compact(s)
	return len(s) != initialLen
}
