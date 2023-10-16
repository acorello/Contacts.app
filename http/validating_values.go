package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func (my UrlValues) IntOrPanic(name string, defaultIfBlank int) int {
	v := my.Trim(name)
	if v == "" {
		return defaultIfBlank
	}
	if i, err := strconv.Atoi(v); err != nil {
		panic(fmt.Errorf("failed to parse %q: %v", name, err))
	} else {
		return i
	}
}

func (my UrlValues) Trim_NotBlank(name string) (string, error) {
	v := my.Trim(name)
	if v == "" {
		return "", fmt.Errorf("%q is blank or empty", name)
	}
	return v, nil
}
