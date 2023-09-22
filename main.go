package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//go:embed *.html
var templates embed.FS

var concactsTemplate = template.Must(template.ParseFS(templates, "contacts.html"))
var newContactTemplate = template.Must(template.ParseFS(templates, "contacts_new.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/contacts/new",
		LoggingHandler(http.HandlerFunc(newContactHandler)))

	mux.HandleFunc("/contacts/",
		LoggingHandler(http.HandlerFunc(contactsHandler)))

	mux.HandleFunc("/",
		LoggingHandler(http.RedirectHandler("/contacts/", http.StatusFound)))

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
}

func newContactHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		getNewContactForm(w, r)
	case http.MethodPost:
		postNewContactForm(w, r)
	default:
		respondErrMethodNotImplemented(w, r)
	}
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		getContacts(w, r)
	default:
		respondErrMethodNotImplemented(w, r)
	}
}

func respondErrMethodNotImplemented(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("method %q not implemented for %q", r.Method, r.URL)
	http.Error(w, msg, http.StatusNotImplemented)
}

func getNewContactForm(w http.ResponseWriter, r *http.Request) {
	templateId := "contacts_new.html"
	args := NewContactForm{Errors: make(ErrorMap)}
	if err := newContactTemplate.ExecuteTemplate(w, templateId, args); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

type NewContactForm struct {
	Contact
	Errors error
}

func postNewContactForm(w http.ResponseWriter, r *http.Request) {
	newContact, err := makeNewContact(r)
	if err != nil {
		log.Printf("%#v", err)
		templateId := "contacts_new.html"
		args := NewContactForm{
			Contact: newContact,
			Errors:  err,
		}
		if err := newContactTemplate.ExecuteTemplate(w, templateId, args); err != nil {
			log.Printf("error rendering template: %v", err)
		}
	} else {
		contactRepository.Store(newContact)
		log.Printf("Stored: %#v", newContact)
		http.Redirect(w, r, "/contacts/", http.StatusFound)
	}
}

type ValidatingValues struct {
	url.Values
	ErrorMap
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

func (my ValidatingValues) NotEmptyString(name string) string {
	v := my.Get(name)
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		my.ErrorMap[name] = "blank or empty"
	}
	return v
}

func makeNewContact(r *http.Request) (c Contact, err error) {
	r.ParseForm()
	form := ValidatingValues{
		Values:   r.Form,
		ErrorMap: make(ErrorMap),
	}
	c.FirstName = form.NotEmptyString("FirstName")
	c.LastName = form.NotEmptyString("LastName")
	c.Email = form.NotEmptyString("Email")
	c.Phone = form.NotEmptyString("Phone")
	if form.HasErrors() {
		return c, form.ErrorMap
	} else {
		return c, nil
	}
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	searchTerm := q.Get("SearchTerm")
	searchTerm = strings.TrimSpace(searchTerm)
	var contacts []Contact
	if searchTerm == "" {
		log.Printf("Listing all contacts")
		contacts = contactRepository.FindAll()
	} else {
		log.Printf("Listing contacts containing %q", searchTerm)
		contacts = contactRepository.FindBySearchTerm(searchTerm)
	}
	templateId := "contacts.html"
	args := map[string]any{
		"SearchTerm": searchTerm,
		"Contacts":   contacts,
	}
	if err := concactsTemplate.ExecuteTemplate(w, templateId, args); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`serving '%s %s'`, r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}

var contactRepository ContactRepository = []Contact{
	{
		Id:        "0",
		FirstName: "Joe",
		LastName:  "Bloggs",
		Phone:     "+44(0)751123456",
		Email:     "joebloggs@example.com",
	},
}

type ContactRepository []Contact

func (me ContactRepository) FindAll() (result []Contact) {
	for _, c := range me {
		result = append(result, c)
	}
	return
}

func (me *ContactRepository) Store(c Contact) error {
	*me = append(*me, c)
	return nil
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
	Id, FirstName, LastName, Phone, Email string
}

func (my Contact) AnyFieldContains(s string) bool {
	p := strings.Contains
	return p(my.FirstName, s) || p(my.LastName, s) || p(my.Phone, s) || p(my.Email, s)
}
