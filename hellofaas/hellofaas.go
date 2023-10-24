// Package helloworld provides a set of Cloud Functions samples.
package hellofaas

import (
	"fmt"
	"html"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloFaaS", HelloFaaS)
}

// HelloFaaS is an HTTP Cloud Function with a request parameter.
func HelloFaaS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.URL.Query()
	if name := q.Get("name"); name == "" {
		fmt.Fprint(w, "Hello, FaaS!")
		return
	} else {
		fmt.Fprintf(w, "Hello, %s-FaaS!", html.EscapeString(name))
	}
}
