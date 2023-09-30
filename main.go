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

	const BASE_PATH = "/contact/"
	contactHandler := http_contact.NewContactHandler(BASE_PATH, &repo)

	mux.HandleFunc(BASE_PATH, LoggingHandler(contactHandler))

	mux.HandleFunc("/static/",
		LoggingHandler(static.FileServer()))

	mux.HandleFunc("/",
		LoggingHandler(http.RedirectHandler("/contact/list", http.StatusFound)))

	address := "localhost:8080"
	log.Printf("Starting server at %q", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`serving '%s %s'`, r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}
