package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"html/template"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/contact/http/ht"

	_http "dev.acorello.it/go/contacts/http"
)

type contactHTTPHandler struct {
	validResourcePaths
	contactRepository contact.Repository
}

func NewContactHandler(paths validResourcePaths, repo contact.Repository) contactHTTPHandler {
	return contactHTTPHandler{
		validResourcePaths: paths,
		contactRepository:  repo,
	}
}

func (h contactHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case h.Root:
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	case h.Form:
		switch r.Method {
		case http.MethodGet:
			h.GetForm(w, r)
		case http.MethodPost:
			h.PostForm(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	case h.List:
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
	id, err := contact.ParseId(q.Get(CustomerId))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", q.Get(CustomerId), err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	contact, found := h.contactRepository.FindById(id)
	hasIdArg := q.Has(CustomerId)
	if hasIdArg && !found {
		w.WriteHeader(http.StatusNotFound)
	} else if !hasIdArg {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// should I move error handling within the template package? maybe better now to just panic?
		urls := ht.ContactPageURLs{
			ContactList: template.URL(h.List),
			ContactForm: h.contactFormURL(contact),
		}
		if err := ht.WriteContact(w, contact, urls); err != nil {
			log.Printf("error rendering template: %v", err)
		}
	}
}

func (h contactHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	form := _http.NewUrlValues(r)
	if !form.Has(CustomerId) {
		msg := fmt.Sprintf("Missing %q from submitted form: %#v", CustomerId, r.Form)
		http.Error(w, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}
	_id := form.Trim(CustomerId)
	id, err := contact.ParseId(_id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", _id, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	h.contactRepository.Delete(id)
	http.Redirect(w, r, h.List, http.StatusSeeOther)
}

func (h contactHTTPHandler) PostForm(w http.ResponseWriter, r *http.Request) {
	var renderingError error
	contact, err := parseContact(r)
	if err != nil {
		log.Printf("Error parsing contacto form: %+v", err)
		contactForm := ht.NewFormWith(contact)
		renderingError = ht.WriteContactForm(w, contactForm, ht.ContactFormPageURLs{
			ContactForm: h.contactFormURL(contact),
		})
	} else {
		h.contactRepository.Store(contact)
		log.Printf("Stored: %#v", contact)
		http.Redirect(w, r, h.List, http.StatusFound)
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func (h contactHTTPHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var renderingError error
	if existingContact := q.Has(CustomerId); !existingContact {
		// blank form to create a new contact
		contactForm := ht.NewForm()
		urls := ht.ContactFormPageURLs{
			ContactForm: template.URL(h.Form),
			ContactList: template.URL(h.List),
		}
		renderingError = ht.WriteContactForm(w, contactForm, urls)
	} else {
		_id := q.Get(CustomerId)
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
			urls := ht.ContactFormPageURLs{
				ContactList:   template.URL(h.List),
				ContactForm:   h.contactFormURL(c),
				DeleteContact: h.contactURL(c),
			}
			renderingError = ht.WriteContactForm(w, ht.NewFormWith(c), urls)
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
	if err := ht.WriteContactList(w, ht.SearchPage{
		SearchTerm: searchTerm,
		Contacts:   contacts,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

func parseContact(r *http.Request) (c contact.Contact, err error) {
	form := _http.NewUrlValues(r)
	errors := make(map[string]error)
	getAndCollect := func(get func(string) (string, error), key string, variable *string) {
		val, err := get(key)
		if err != nil {
			errors[key] = err
		}
		*variable = val
	}
	if _id := form.Trim(CustomerId); _id == "" {
		c.Id = contact.NewId()
		// TODO: prhaps I shall differentiate this case using PUT / POST
		log.Printf("Got blank contact id assuming new contact, assigning new id %q", c.Id)
	} else if id, err := contact.ParseId(_id); err != nil {
		errors[CustomerId] = err
	} else {
		c.Id = id
	}
	getAndCollect(form.Trim_NotBlank, "FirstName", &c.FirstName)
	getAndCollect(form.Trim_NotBlank, "LastName", &c.LastName)
	getAndCollect(form.Trim_NotBlank, "Email", &c.Email)
	c.Phone = form.Trim("Phone")
	if len(errors) > 0 {
		return c, fmt.Errorf("%#v", errors)
	} else {
		return c, nil
	}
}
