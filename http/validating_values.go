package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type UrlValues struct {
	url.Values
}

func NewUrlValues(r *http.Request) UrlValues {
	r.ParseForm()
	return UrlValues{
		Values: r.Form,
	}
}

func (my UrlValues) Trim(name string) string {
	v := my.Get(name)
	v = strings.TrimSpace(v)
	return v
}

func (my UrlValues) Trim_NotBlank(name string) (string, error) {
	v := my.Get(name)
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		return "", fmt.Errorf("%q is blank or empty", name)
	}
	return v, nil
}
