package ht

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/seq"
	"dev.acorello.it/go/contacts/templates"
)

//go:embed *.html
var myTemplates embed.FS

var contactTemplate,
	contactFormTemplate,
	contactsTemplate *template.Template

func init() {
	contactTemplate = makeTemplate(myTemplates, "contact.html")
	contactsTemplate = makeTemplate(myTemplates, "contacts.html")
	contactFormTemplate = makeTemplate(myTemplates, "contact_form.html")
}

func makeTemplate(files fs.FS, templateFile string) *template.Template {
	t := template.Must(template.ParseFS(templates.CommonFS(), "layout.html"))
	template.Must(t.ParseFS(files, templateFile))
	names := seq.Map(t.Templates(), (*template.Template).Name)
	log.Printf("%q associated templates: %#v", templateFile, names)
	return t
}

type ContactPageURLs struct {
	ContactList, ContactForm template.URL
}

func WriteContact(w io.Writer, c contact.Contact, u ContactPageURLs) error {
	return contactTemplate.Execute(w, map[string]any{
		"Contact": c,
		"URLs":    u,
	})
}

type ContactForm struct {
	contact.Contact
	Errors templates.ErrorMap
}

func NewFormWith(c contact.Contact) ContactForm {
	return ContactForm{
		Contact: c,
		Errors:  templates.NewErrorMap(),
	}
}

func NewForm() ContactForm {
	return ContactForm{
		Errors: templates.NewErrorMap(),
	}
}

func WriteContactForm(w io.Writer, c ContactForm) error {
	return contactFormTemplate.Execute(w, c)
}

type SearchPage struct {
	SearchTerm string
	Contacts   []contact.Contact
}

func WriteContacts(w io.Writer, s SearchPage) error {
	return contactsTemplate.Execute(w, s)
}
