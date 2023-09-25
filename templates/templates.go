package templates

import (
	"embed"
	"fmt"
)

//go:embed *.html
var fs embed.FS

func CommonFS() embed.FS {
	return fs
}

type ErrorMap map[string]error

func NewErrorMap() ErrorMap {
	return make(ErrorMap)
}

func (my ErrorMap) Error() string {
	return fmt.Sprintf("%#v", my)
}
