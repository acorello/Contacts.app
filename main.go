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
	searchTerm := q.Get("SearchTerm")
	searchTerm = strings.TrimSpace(searchTerm)
	var contacts []Contact
	if searchTerm == "" {
		contacts = contactRepository.FindAll()
	} else {
		contacts = contactRepository.FindBySearchTerm(searchTerm)
	}
	err := concactsTemplate.ExecuteTemplate(w, "contacts.html",
		map[string]any{
			"SearchTerm": searchTerm,
			"Contacts":   contacts,
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

var contactRepository ContactRepository = []Contact{
	{
		Id:    "0",
		First: "Joe",
		Last:  "Bloggs",
		Phone: "+44(0)751123456",
		Email: "joebloggs@example.com",
	},
}

type ContactRepository []Contact

func (me ContactRepository) FindAll() (result []Contact) {
	for _, c := range me {
		result = append(result, c)
	}
	return
}

func (me ContactRepository) FindBySearchTerm(term string) (result []Contact) {
	for _, c := range me {
		if c.AnyFieldContains(term) {
			result = append(result, c)
		}
	}
	return
}

type Contact struct {
	Id, First, Last, Phone, Email string
}

func (my Contact) AnyFieldContains(s string) bool {
	p := strings.Contains
	return p(my.First, s) || p(my.Last, s) || p(my.Phone, s) || p(my.Email, s)
}
