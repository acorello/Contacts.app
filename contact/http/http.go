package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/contact/template"
	"dev.acorello.it/go/contacts/http_util"
	"github.com/google/uuid"
)

// I want to have a single handler for the /contact path and subpaths
// for instance /contact/form
// Why do I want a /contact/form?
// Why not using a POST to /contact
// POST /contact may work semantically, but when I want to create a new contact, I have to GET an empty form
// GET /contact?Id=<ID> is already taken by the
// GET /contact?new could work if I dispatch on the params… however params are a map of value, they allow ambiguous requests like /contact?new&Id=<ID>
// GET /contact/form instead is unique.
// But in order for this to work we have to make this module path-aware… which actually it's OK, because this module is responsible of handling the HTTP requests for a resource, so has to know about parameter names and format, for example.
//

const (
	BASE_PATH       = "/contact/"
	form_path       = "/contact/form"
	collection_path = "/contact/list"
)

type contactHTTPHandler struct {
	contactRepository contact.Repository
}

func NewContactHandler(r contact.Repository) contactHTTPHandler {
	return contactHTTPHandler{
		contactRepository: r,
	}
}

func (h contactHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case BASE_PATH:
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http_util.RespondErrMethodNotImplemented(w, r)
		}
	case form_path:
		switch r.Method {
		case http.MethodGet:
			h.GetForm(w, r)
		case http.MethodPost:
			h.PostForm(w, r)
		default:
			http_util.RespondErrMethodNotImplemented(w, r)
		}
	case collection_path:
		switch r.Method {
		case http.MethodGet:
			h.GetList(w, r)
		default:
			http_util.RespondErrMethodNotImplemented(w, r)
		}
	default:
		errorMsg := fmt.Sprintf("Path %q is not supported", p)
		http.Error(w, errorMsg, http.StatusNotFound)
	}
}

func (h contactHTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("Id")
	id = strings.TrimSpace(id)
	contact, found := h.contactRepository.FindById(id)
	hasIdArg := q.Has("Id")
	if hasIdArg && !found {
		w.WriteHeader(http.StatusNotFound)
	} else if !hasIdArg {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		template.WriteContactHTML(w, contact)
	}
}

func (h contactHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	form := http_util.NewUrlValues(r)
	var renderingError error
	if !form.Has("Id") {
		msg := fmt.Sprintf("Missing %q from submitted form: %#v", "Id", r.Form)
		http.Error(w, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}
	id := form.Trim("Id") //TODO: rename additional readers with Get prefix
	h.contactRepository.Delete(id)
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func (h contactHTTPHandler) PostForm(w http.ResponseWriter, r *http.Request) {
	var renderingError error
	contactForm, err := parseContactForm(r)
	if err != nil {
		log.Printf("%#v", err)
		contactForm := template.NewFormWith(contactForm)
		renderingError = template.WriteContactFormHTML(w, contactForm)
	} else {
		if contactForm.Id == "" {
			contactForm.Id = uuid.NewString()
		}
		h.contactRepository.Store(contactForm)
		log.Printf("Stored: %#v", contactForm)
		http.Redirect(w, r, "/contacts", http.StatusFound)
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func (h contactHTTPHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var renderingError error
	if editContact := q.Has("Id"); !editContact {
		// blank form to create a new contact
		contactForm := template.NewForm()
		renderingError = template.WriteContactFormHTML(w, contactForm)
	} else {
		id := q.Get("Id")
		id = strings.TrimSpace(id)
		c, found := h.contactRepository.FindById(id)
		if !found {
			w.WriteHeader(http.StatusNotFound)
		} else {
			renderingError = template.WriteContactFormHTML(w, template.NewFormWith(c))
		}
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func (h contactHTTPHandler) GetList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	searchTerm := q.Get("SearchTerm")
	searchTerm = strings.TrimSpace(searchTerm)
	var contacts []contact.Contact
	if searchTerm == "" {
		log.Printf("Listing all contacts")
		contacts = h.contactRepository.FindAll()
	} else {
		log.Printf("Listing contacts containing %q", searchTerm)
		contacts = h.contactRepository.FindBySearchTerm(searchTerm)
	}
	if err := template.WriteContactsHTML(w, template.SearchPage{
		SearchTerm: searchTerm,
		Contacts:   contacts,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

func parseContactForm(r *http.Request) (c contact.Contact, err error) {
	form := http_util.NewUrlValues(r)
	errors := map[string]error{}
	getAndCollect := func(get func(string) (string, error), key string, store *string) {
		val, err := get(key)
		if err != nil {
			errors[key] = err
		}
		*store = val
	}
	c.Id = form.Trim("Id")
	getAndCollect(form.Trim_NotBlank, "FirstName", &c.FirstName)
	getAndCollect(form.Trim_NotBlank, "LastName", &c.LastName)
	getAndCollect(form.Trim_NotBlank, "Email", &c.Email)
	getAndCollect(form.Trim_NotBlank, "Phone", &c.Phone)
	if len(errors) > 0 {
		return c, fmt.Errorf("%#v", errors)
	} else {
		return c, nil
	}
}
