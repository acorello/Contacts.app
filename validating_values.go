package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ValidatingValues struct {
	url.Values
	ErrorMap
}

func NewValidatingValues(r *http.Request) ValidatingValues {
	r.ParseForm()
	return ValidatingValues{
		Values:   r.Form,
		ErrorMap: make(ErrorMap),
	}
}

type ErrorMap map[string]string

func (my ErrorMap) Error() string {
	return fmt.Sprintf("%#v", my)
}

func (my ValidatingValues) HasErrors() bool {
	return len(my.ErrorMap) > 0
}

func (my ValidatingValues) ErrorsMap() map[string]string {
	return my.ErrorMap
}

func (my ValidatingValues) String(name string) string {
	v := my.Get(name)
	v = strings.TrimSpace(v)
	return v
}

func (my ValidatingValues) NotEmptyString(name string) string {
	v := my.Get(name)
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		my.ErrorMap[name] = "blank or empty"
	}
	return v
}
