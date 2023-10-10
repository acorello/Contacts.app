package http

import (
	"cmp"
	"fmt"
	"html/template"
	"net/url"
	"slices"

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

func (my validResourcePaths) contactFormURL(c contact.Contact) template.URL {
	q := url.Values{}
	q.Add("Id", c.Id.String())
	u := url.URL{
		Path:     my.Form,
		RawQuery: q.Encode(),
	}
	return template.URL(u.String())
}

func (my validResourcePaths) patchContactEmailURL(c contact.Contact) template.URL {
	q := url.Values{}
	q.Add("Id", c.Id.String())
	u := url.URL{
		Path:     my.Email,
		RawQuery: q.Encode(),
	}
	return template.URL(u.String())
}

func (my validResourcePaths) contactURL(c contact.Contact) template.URL {
	q := url.Values{}
	q.Add("Id", c.Id.String())
	u := url.URL{
		Path:     my.Root,
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
