package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/contact/template"
	_http "dev.acorello.it/go/contacts/http"
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

type contactHTTPHandler struct {
	basePath string
	formPath string
	ListPath string

	contactRepository contact.Repository
}

// basePath should be absolute, end with '/', and have at least one element
func NewContactHandler(basePath string, r contact.Repository) contactHTTPHandler {
	const docMsg = "path should be absolute, end with '/', and have at least one element"
	if !(len(basePath) > 2 && basePath[0] == '/' && basePath[len(basePath)] == '/') {
		// panic: contract was violated, server requires this initialization
		panic(fmt.Sprintf("%s, got %q", docMsg, basePath))
	}
	return contactHTTPHandler{
		basePath: basePath,
		formPath: basePath + "form",
		ListPath: basePath + "list",

		contactRepository: r,
	}
}

// expects to be bound to BASE_PATH, a folder. Will dispatch on any sub-path.
func (h contactHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case h.basePath:
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	case h.formPath:
		switch r.Method {
		case http.MethodGet:
			h.GetForm(w, r)
		case http.MethodPost:
			h.PostForm(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	case h.ListPath:
		switch r.Method {
		case http.MethodGet:
			h.GetList(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	default:
		errorMsg := fmt.Sprintf("Path %q is not supported", p)
		http.Error(w, errorMsg, http.StatusNotFound)
	}
}

func (h contactHTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id, err := contact.ParseId(q.Get("Id"))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", q.Get("Id"), err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
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
	form := _http.NewUrlValues(r)
	var renderingError error
	if !form.Has("Id") {
		msg := fmt.Sprintf("Missing %q from submitted form: %#v", "Id", r.Form)
		http.Error(w, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}
	_id := form.Trim("Id")
	id, err := contact.ParseId(_id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", _id, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	h.contactRepository.Delete(id)
	http.Redirect(w, r, h.ListPath, http.StatusSeeOther)
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
			contactForm.Id = contact.NewId()
		}
		h.contactRepository.Store(contactForm)
		log.Printf("Stored: %#v", contactForm)
		http.Redirect(w, r, h.ListPath, http.StatusFound)
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
		_id := q.Get("Id")
		id, err := contact.ParseId(_id)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to parse id %q: %v", _id, err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
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
	form := _http.NewUrlValues(r)
	errors := map[string]error{}
	getAndCollect := func(get func(string) (string, error), key string, store *string) {
		val, err := get(key)
		if err != nil {
			errors[key] = err
		}
		*store = val
	}
	id, err := contact.ParseId(form.Trim("Id"))
	if err != nil {
		errors["Id"] = err
	}
	c.Id = id
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
