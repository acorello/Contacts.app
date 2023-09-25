package template

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

func WriteContactHTML(w io.Writer, c contact.Contact) error {
	return contactTemplate.Execute(w, c)
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

func WriteContactFormHTML(w io.Writer, c ContactForm) error {
	return contactFormTemplate.Execute(w, c)
}

type SearchPage struct {
	SearchTerm string
	Contacts   []contact.Contact
}

func WriteContactsHTML(w io.Writer, s SearchPage) error {
	return contactsTemplate.Execute(w, s)
}
