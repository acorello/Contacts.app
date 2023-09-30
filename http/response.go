package http

import (
	"fmt"
	"net/http"
)

func RespondErrMethodNotImplemented(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("method %q not implemented for %q", r.Method, r.URL)
	http.Error(w, msg, http.StatusNotImplemented)
}
