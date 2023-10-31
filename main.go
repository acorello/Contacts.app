package main

import (
	"log"
	"net/http"
	"os"

	"dev.acorello.it/go/contacts/contact"
	http_contact "dev.acorello.it/go/contacts/contact/http"
	"dev.acorello.it/go/contacts/public_assets"
)

var repo = contact.NewPopulatedInMemoryContactRepository()

func main() {
	mux := http.NewServeMux()
	const publicRootPath = "/public/"
	mux.HandleFunc(publicRootPath,
		LoggingHandler(http.StripPrefix(publicRootPath, public_assets.FileServer())))

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

	address := bindAddress()
	log.Printf("Starting server at %q", address)
	if serverErr := http.ListenAndServe(address, mux); serverErr != nil {
		log.Fatal(serverErr)
	}
}

func bindAddress() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	return host + ":8080"
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`serving '%s %s'`, r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}
