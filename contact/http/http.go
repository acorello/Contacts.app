package http

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

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
	case h.Email:
		switch r.Method {
		case http.MethodPatch:
			h.PatchEmail(w, r)
		default:
			_http.RespondErrMethodNotImplemented(w, r)
		}
	default:
		errorMsg := fmt.Sprintf("Path %q is not supported", p)
		log.Println(errorMsg)
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
			ContactForm: h.contactResourceURL(contact, h.Form),
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
	contact, errors := parseContact(r)
	if errors != nil && len(errors) > 0 {
		log.Printf("Error parsing contact form: %+v", errors)
		contactForm := ht.NewFormWith(contact)
		contactForm.Errors = errors
		renderingError = ht.WriteContactForm(w, contactForm, ht.ContactFormPageURLs{
			ContactForm:       h.contactResourceURL(contact, h.Form),
			PatchContactEmail: h.contactResourceURL(contact, h.Email),
		})
	} else {
		// TODO: implement validation (eg. [e-mail]--N--1--[contactId] ) and error handling
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
			ContactForm:       template.URL(h.Form),
			ContactList:       template.URL(h.List),
			PatchContactEmail: h.contactResourceURL(contactForm.Contact, h.Email),
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
		contact, found := h.contactRepository.FindById(id)
		if !found {
			w.WriteHeader(http.StatusNotFound)
		} else {
			urls := ht.ContactFormPageURLs{
				ContactList:       template.URL(h.List),
				ContactForm:       h.contactResourceURL(contact, h.Form),
				DeleteContact:     h.contactResourceURL(contact, h.Root),
				PatchContactEmail: h.contactResourceURL(contact, h.Email),
			}
			renderingError = ht.WriteContactForm(w, ht.NewFormWith(contact), urls)
		}
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

func (h contactHTTPHandler) PatchEmail(w http.ResponseWriter, r *http.Request) {
	q := _http.NewUrlValues(r)
	contactId := contact.Id(q.Trim("Id"))
	contactEmail := q.Trim("Email")
	log.Printf("validating e-mail %q for contactId %q", contactEmail, contactId)
	existingContactId, found := h.contactRepository.FindIdByEmail(contactEmail)
	if found && existingContactId != contactId {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "email address already in use")
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h contactHTTPHandler) GetList(w http.ResponseWriter, r *http.Request) {
	q := _http.NewUrlValues(r)
	searchTerm := q.Trim("SearchTerm")
	page := contact.Page{
		Offset: q.IntOrPanic("pageOffset", 0),
		Size:   q.IntOrPanic("pageSize", 0),
	}
	page.Offset = max(page.Offset, 0)
	page.Size = max(page.Size, 10)
	page.Size = min(page.Size, 50)

	var contacts []contact.Contact
	var more bool
	if searchTerm == "" {
		log.Printf("Listing all contacts")
		contacts, more = h.contactRepository.FindAll(page)
	} else {
		log.Printf("Listing contacts containing %q", searchTerm)
		contacts, more = h.contactRepository.FindBySearchTerm(searchTerm, page)
	}
	var nextPageURL template.URL
	if more {
		nextPageURL = h.searchPageURL(page.Next(), searchTerm)
	}
	templateParams := ht.SearchPage{
		SearchTerm: searchTerm,
		Contacts:   contacts,
		URLs: ht.SearchPageURLs{
			NextPage: nextPageURL,
		},
	}
	if err := ht.WriteContactList(w, templateParams); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

var nameRegEx = re{regexp.MustCompile(`^\w+(?:[- ']\w+)*$`)}

func parseContact(r *http.Request) (c contact.Contact, err map[string]error) {
	form := _http.NewUrlValues(r)
	errors := make(map[string]error)
	getAndCollect := func(get func(string) (string, error), key string, variable *string, validators ...func(string) error) {
		val, err := get(key)
		*variable = val
		if err != nil {
			errors[key] = err
			return
		}
		for _, v := range validators {
			err := v(val)
			if err != nil {
				errors[key] = err
				return
			}
		}
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

	getAndCollect(form.Trim_NotBlank, "FirstName", &c.FirstName, nameRegEx.Validate)
	getAndCollect(form.Trim_NotBlank, "LastName", &c.LastName)
	getAndCollect(form.Trim_NotBlank, "Email", &c.Email)
	c.Phone = form.Trim("Phone")
	return c, errors
}

type re struct {
	*regexp.Regexp
}

func (me re) Validate(val string) error {
	if !me.MatchString(val) {
		return fmt.Errorf("%q did not match %#q", val, me.Regexp.String())
	}
	return nil
}
