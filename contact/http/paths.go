package http

import (
	"html/template"
	"net/url"
	"strings"
)

const (
	CustomerId = "Id"
)

type Path string

func (me Path) Add(param, value string) Path {
	params := url.Values{}
	params.Add(param, value)
	_current := string(me)
	var separator string
	if strings.Contains(_current, "?") {
		separator = "&"
	} else {
		separator = "?"
	}
	return Path(_current + separator + params.Encode())
}

func (me Path) TemplateURL() template.URL {
	return template.URL(me)
}

func (me Path) String() string {
	return string(me)
}
