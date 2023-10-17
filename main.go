package main

import (
	"log"
	"net/http"

	"dev.acorello.it/go/contacts/contact"
	http_contact "dev.acorello.it/go/contacts/contact/http"
	"dev.acorello.it/go/contacts/static"
)

var repo = contact.NewPopulatedInMemoryContactRepository()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/static/", LoggingHandler(static.FileServer()))

	contactResourcePaths := http_contact.ResourcePaths{
		Root:  "/contact/",
		Form:  "/contact/form",
		List:  "/contact/list",
		Email: "/contact/email",
	}

	if validatedPaths, err := contactResourcePaths.Validated(); err != nil {
		log.Fatal(err)
	} else {
		contactHandler := http_contact.NewContactHandler(validatedPaths, &repo)
		mux.HandleFunc(validatedPaths.Root.String(), LoggingHandler(contactHandler))
		homeRedirect := http.RedirectHandler(validatedPaths.List.String(), http.StatusFound)
		mux.HandleFunc("/", LoggingHandler(homeRedirect))
	}

	address := "localhost:8080"
	log.Printf("Starting server at %q", address)
	if serverErr := http.ListenAndServe(address, mux); serverErr != nil {
		log.Fatal(serverErr)
	}
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`serving '%s %s'`, r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}
