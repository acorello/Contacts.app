package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/contact/http/ht"
	"dev.acorello.it/go/contacts/seq"
	"github.com/acorello/must"
	"github.com/acorello/uttpil"
)

type Paths struct {
	Root, Form, List, Email Path
}

type paths Paths

// Validated checks that:
// - paths are distinct
func (my Paths) Validated() (v paths, err error) {
	if seq.HasDuplicates(my.Root, my.Form, my.List, my.Email) {
		return v, fmt.Errorf("path elements must be unique. Got %+v", my)
	}
	return paths(my), nil
}

func RegisterHandlers(mux *http.ServeMux, paths paths, repo contact.Repository) {
	h := contactHTTPHandler{
		paths:             paths,
		contactRepository: repo,
	}
	mux.Handle(paths.Root.String(), uttpil.ForMethod{
		GET:    h.Get,
		DELETE: h.Delete,
	})
	mux.Handle(paths.Form.String(), uttpil.ForMethod{
		GET:  h.GetForm,
		POST: h.PostForm,
	})
	mux.Handle(paths.List.String(), uttpil.ForMethod{
		GET: h.GetList,
	})
	mux.Handle(paths.Email.String(), uttpil.ForMethod{
		PATCH: h.PatchEmail,
	})
}

type contactHTTPHandler struct {
	paths             paths
	contactRepository contact.Repository
}

func (h contactHTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has(CustomerId) {
		w.WriteHeader(http.StatusBadRequest)
	} else if id, err := contact.ParseId(q.Get(CustomerId)); err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", q.Get(CustomerId), err)
		http.Error(w, errMsg, http.StatusBadRequest)
	} else if theContact, found := h.contactRepository.FindById(id); !found {
		w.WriteHeader(http.StatusNotFound)
	} else {
		_id := theContact.Id.String()
		urls := ht.ContactPageURLs{
			ContactList: template.URL(h.paths.List),
			ContactForm: h.paths.Form.Add(CustomerId, _id).TemplateURL(),
		}
		err = ht.WriteContact(w, theContact, urls)
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
	}
}

func (h contactHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	form, err := uttpil.NewUrlValuesHelper(r)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	if !form.Has(CustomerId) {
		msg := fmt.Sprintf("Missing %q from submitted form: %#v", CustomerId, r.Form)
		http.Error(w, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}
	_id := form.Get(CustomerId, strings.TrimSpace)
	id, err := contact.ParseId(_id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse id %q: %v", _id, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	h.contactRepository.Delete(id)
	http.Redirect(w, r, h.paths.List.String(), http.StatusSeeOther)
}

func (h contactHTTPHandler) PostForm(w http.ResponseWriter, r *http.Request) {
	var renderingError error
	form, err := uttpil.NewUrlValuesHelper(r) // confusingly enough, the library decodes a form into url.Values,
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	theContact, errors := parseContact(form)
	if errors != nil && len(errors) > 0 {
		log.Printf("Error parsing contact form: %+v", errors)
		contactForm := ht.NewFormWith(theContact)
		contactForm.Errors = errors
		_id := theContact.Id.String()
		renderingError = ht.WriteContactForm(w, ht.ContactFormPage{
			ContactForm: contactForm,
			URLs: ht.ContactFormPageURLs{
				ContactForm:       h.paths.Form.Add(CustomerId, _id).TemplateURL(),
				PatchContactEmail: h.paths.Email.Add(CustomerId, _id).TemplateURL(),
			},
		})
	} else {
		h.contactRepository.Store(theContact)
		log.Printf("Stored: %#v", theContact)
		http.Redirect(w, r, h.paths.List.String(), http.StatusFound)
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
		contactForm.Id = contact.NewId()
		_id := contactForm.Id.String()
		urls := ht.ContactFormPageURLs{
			ContactForm:       h.paths.Form.Add(CustomerId, _id).TemplateURL(),
			ContactList:       h.paths.List.TemplateURL(),
			PatchContactEmail: h.paths.Email.Add(CustomerId, _id).TemplateURL(),
		}
		renderingError = ht.WriteContactForm(w, ht.ContactFormPage{
			ContactForm: contactForm,
			URLs:        urls,
		})
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
			_id := contact.Id.String()
			urls := ht.ContactFormPageURLs{
				ContactList:       h.paths.List.TemplateURL(),
				ContactForm:       h.paths.Form.Add(CustomerId, _id).TemplateURL(),
				DeleteContact:     h.paths.Root.Add(CustomerId, _id).TemplateURL(),
				PatchContactEmail: h.paths.Email.Add(CustomerId, _id).TemplateURL(),
			}
			renderingError = ht.WriteContactForm(w, ht.ContactFormPage{
				ContactForm: ht.NewFormWith(contact),
				URLs:        urls,
			})
		}
	}
	if renderingError != nil {
		log.Printf("error rendering template: %v", renderingError)
	}
}

// PatchEmail is only used for validation purposes at the moment
func (h contactHTTPHandler) PatchEmail(w http.ResponseWriter, r *http.Request) {
	q, err := uttpil.NewUrlValuesHelper(r)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	contactId := contact.Id(q.Get(CustomerId, strings.TrimSpace))
	contactEmail := q.Get("Email", strings.TrimSpace)
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
	q, err := uttpil.NewUrlValuesHelper(r)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	searchTerm := q.Get("SearchTerm", strings.TrimSpace)
	// does it make sense to create a separate object dedicated to the parsing?
	page := contact.Page{
		Offset: must.Get(asInt(q.Get("pageOffset"), 0)),
		Size:   must.Get(asInt(q.Get("pageSize"), 0)),
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
		nextPageURL = searchPageURL(page.Next(), searchTerm, h.paths.List.String())
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

func asInt(s string, whenBlank int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return whenBlank, nil
	}
	if v, err := strconv.Atoi(s); err != nil {
		return 0, fmt.Errorf("failed to parse %q: %v", s, err)
	} else {
		return v, nil
	}
}

func searchPageURL(page contact.Page, searchTerm, searchPagePath string) template.URL {
	q := url.Values{}
	if searchTerm != "" {
		q.Add("SearchTerm", searchTerm)
	}
	q.Add("pageOffset", strconv.Itoa(page.Offset))
	q.Add("pageSize", strconv.Itoa(page.Size))
	u := url.URL{
		Path:     searchPagePath,
		RawQuery: q.Encode(),
	}
	return template.URL(u.String())
}

var nameRegEx = re{regexp.MustCompile(`^\w+(?:[- ']\w+)*$`)}

func parseContact(form uttpil.UrlValuesHelper) (c contact.Contact, err map[string]error) {
	form.Give(CustomerId, func(val string) error {
		val = strings.TrimSpace(val)
		if val == "" {
			c.Id = contact.NewId()
			log.Printf("Got blank contact id assuming new contact, assigning new id %q", c.Id)
			return nil
		}
		if id, err := contact.ParseId(val); err != nil {
			return err
		} else {
			c.Id = id
			return nil
		}
	})
	form.Give("FirstName", func(value string) error {
		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("blank")
		} else if !nameRegEx.MatchString(value) {
			return fmt.Errorf("invalid name")
		}
		c.FirstName = value
		return nil
	})
	form.Give("LastName", func(value string) error {
		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("blank")
		}
		c.LastName = value
		return nil
	})
	form.Give("Email", func(value string) error {
		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("blank")
		}
		c.Email = value
		return nil
	})
	form.Give("Phone", func(value string) error {
		c.Phone = strings.TrimSpace(value)
		return nil
	})
	return c, form.Errors()
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
