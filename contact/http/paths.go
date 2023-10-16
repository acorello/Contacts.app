package http

import (
	"html/template"
	"net/url"
	"strings"
)

const (
	CustomerId = "Id"
)

type ResourcePath string

func (me ResourcePath) Add(param, value string) ResourcePath {
	params := url.Values{}
	params.Add(param, value)
	_current := string(me)
	var separator string
	if strings.Contains(_current, "?") {
		separator = "&"
	} else {
		separator = "?"
	}
	return ResourcePath(_current + separator + params.Encode())
}

func (me ResourcePath) Path() string {
	return string(me)
}

func (me ResourcePath) TemplateURL() template.URL {
	return template.URL(me)
}

func (me ResourcePath) String() string {
	return string(me)
}
