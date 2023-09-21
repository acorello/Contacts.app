package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed *.html
var templates embed.FS

var concactsTemplate = template.Must(template.ParseFS(templates, "contacts.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/",
		LoggingHandler(http.RedirectHandler("/contacts/", http.StatusFound)))

	mux.HandleFunc("/contacts/",
		LoggingHandler(http.HandlerFunc(contactsHandler)))

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		getContacts(w, r)
	default:
		msg := fmt.Sprintf("method %q not implemented for %q", method, r.URL)
		http.Error(w, msg, http.StatusNotImplemented)
	}
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	keyword := q.Get("q")
	keyword = strings.TrimSpace(keyword)
	var contacts []Contact
	if keyword == "" {
		contacts = contactRepository.FindAll()
	} else {
		contacts = contactRepository.FindByKeyword(keyword)
	}
	err := concactsTemplate.ExecuteTemplate(w, "contacts.html", map[string]any{
		"Contacts": contacts,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving %q", r.URL)
		h.ServeHTTP(w, r)
	}
}

var contactRepository ContactRepository = make([]Contact, 0)

type ContactRepository []Contact

func (me ContactRepository) FindAll() (result []Contact) {
	for _, c := range me {
		result = append(result, c)
	}
	return
}

func (me ContactRepository) FindByKeyword(keyword string) (result []Contact) {
	for _, c := range me {
		if strings.Contains(c.Name, keyword) {
			result = append(result, c)
		}
	}
	return
}

type Contact struct {
	Name string
}
