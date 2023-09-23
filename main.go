package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

//go:embed *.html
var templates embed.FS

//go:embed static/*
var static embed.FS

var htmlTemplates = parsedTemplateOrPanic("layout.html", "contacts.html", "contact.html", "contact_form.html")

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/contact/form",
		LoggingHandler(http.HandlerFunc(contactFormHandler)))

	mux.HandleFunc("/contacts",
		LoggingHandler(http.HandlerFunc(contactsHandler)))

	mux.HandleFunc("/contact",
		LoggingHandler(http.HandlerFunc(contactHandler)))

	mux.HandleFunc("/static/",
		LoggingHandler(http.FileServer(http.FS(static))))

	mux.HandleFunc("/",
		LoggingHandler(http.RedirectHandler("/contacts", http.StatusFound)))

	address := "localhost:8080"
	log.Printf("Starting server at %q", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case http.MethodGet:
		getContact(w, r)
	case http.MethodDelete:
		deleteContact(w, r)
	default:
		respondErrMethodNotImplemented(w, r)
	}
}

func contactFormHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		getContactForm(w, r)
	case http.MethodPost:
		postContactForm(w, r)
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

func getContactForm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var renderingError error
	if editContact := q.Has("Id"); !editContact {
		// blank form to create a new contact
		contactForm := ContactForm{Errors: make(ErrorMap)}
		renderingError = htmlTemplates.ExecuteTemplate(w, "contact_form.html", contactForm)
	} else {
		id := q.Get("Id")
		id = strings.TrimSpace(id)
		contact, found := contactRepository.FindById(id)
		if !found {
			w.WriteHeader(http.StatusNotFound)
		} else {
			renderingError = htmlTemplates.ExecuteTemplate(w, "contact_form.html", ContactForm{
				Contact: contact,
				Errors:  make(ErrorMap),
			})
		}
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

type ContactForm struct {
	Contact
	Errors error
}

func postContactForm(w http.ResponseWriter, r *http.Request) {
	form := NewValidatingValues(r)
	var renderingError error
	if form.Has("_DELETE_") {
		if !r.Form.Has("Id") {
			w.WriteHeader(http.StatusBadRequest)
		}
		id := form.Get("Id") //TODO: rename additional readers with Get prefix
		contactRepository.Delete(id)
	}
	contact, err := parseContactForm(r)
	if err != nil {
		log.Printf("%#v", err)
		contactForm := ContactForm{
			Contact: contact,
			Errors:  err,
		}
		renderingError = htmlTemplates.ExecuteTemplate(w, "contact_form.html", contactForm)
	} else {
		if contact.Id == "" {
			contact.Id = uuid.NewString()
		}
		contactRepository.Store(contact)
		log.Printf("Stored: %#v", contact)
		http.Redirect(w, r, "/contacts", http.StatusFound)
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func parseContactForm(r *http.Request) (c Contact, err error) {
	form := NewValidatingValues(r)
	c.Id = form.String("Id")
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

func getContact(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("Id")
	id = strings.TrimSpace(id)
	contact, found := contactRepository.FindById(id)
	hasIdArg := q.Has("Id")
	if hasIdArg && !found {
		w.WriteHeader(http.StatusNotFound)
	} else if !hasIdArg {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		htmlTemplates.ExecuteTemplate(w, "contact.html", contact)
	}
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	form := NewValidatingValues(r)
	var renderingError error
	if !form.Has("Id") {
		msg := fmt.Sprintf("Missing %q from submitted form: %#v", "Id", r.Form)
		http.Error(w, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}
	id := form.String("Id") //TODO: rename additional readers with Get prefix
	contactRepository.Delete(id)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
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
	args := map[string]any{
		"SearchTerm": searchTerm,
		"Contacts":   contacts,
	}
	if err := htmlTemplates.ExecuteTemplate(w, "contacts.html", args); err != nil {
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

func (me ContactRepository) FindById(id string) (c Contact, found bool) {
	for _, c := range me {
		if c.Id == id {
			return c, true
		}
	}
	return c, false
}

func (me *ContactRepository) Delete(id string) {
	contacts := *me
	for i, c := range contacts {
		if c.Id == id {
			contacts[i] = contacts[len(contacts)-1]
			contacts[len(contacts)-1] = Contact{}
			*me = contacts[:len(contacts)-1]
			return
		}
	}
}

func (me ContactRepository) FindAll() (result []Contact) {
	for _, c := range me {
		result = append(result, c)
	}
	return
}

func (me *ContactRepository) Store(c Contact) error {
	for i, x := range *me {
		if x.Id == c.Id {
			(*me)[i] = c
			return nil
		}
	}
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

func parsedTemplateOrPanic(file ...string) *template.Template {
	return template.Must(template.ParseFS(templates, file...))
}
