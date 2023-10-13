package http

import (
	"cmp"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"strconv"

	"dev.acorello.it/go/contacts/contact"
)

const (
	CustomerId = "Id"
)

type ResourcePaths struct {
	Root, Form, List, Email string
}

type validResourcePaths ResourcePaths

// paths should be distinct or this will panic
func (my ResourcePaths) Validated() (v validResourcePaths, err error) {
	if hasDuplicates(my.Root, my.Form, my.List, my.Email) {
		return v, fmt.Errorf("path elements must be unique. Got %+v", my)
	}
	return validResourcePaths(my), nil
}

func contactResourceURL(c contact.Contact, resourcePath string) template.URL {
	q := url.Values{}
	q.Add("Id", c.Id.String())
	u := url.URL{
		Path:     resourcePath,
		RawQuery: q.Encode(),
	}
	return template.URL(u.String())
}

func searchPageURL(page contact.Page, searchTerm, searchPagePath string) template.URL {
	q := url.Values{}
	if searchTerm != "" {
		q.Add("SearchTerm", searchTerm)
	}
	q.Add("pageOffset", strconv.Itoa(page.Offset))
	q.Add("pageSize", strconv.Itoa(page.Size))
	u := url.URL{
		Path:     searchPagePath,
		RawQuery: q.Encode(),
	}
	return template.URL(u.String())
}

func hasDuplicates[T cmp.Ordered](s ...T) bool {
	initialLen := len(s)
	slices.Sort(s)
	s = slices.Compact(s)
	return len(s) != initialLen
}
