package template_test

import (
	"fmt"
	"strings"
	"testing"

	std_template "html/template"

	"dev.acorello.it/go/contacts/contact"
	"dev.acorello.it/go/contacts/contact/template"
	"dev.acorello.it/go/contacts/templates"
	"golang.org/x/net/html"
)

var aContact = contact.Contact{
	Id:        "CNT_1234",
	FirstName: "FIRST_NAME",
	LastName:  "LAST_NAME",
	Phone:     "PHONE",
	Email:     "EMAIL",
}

// Check HTML is syntactically valid and that it contains all properties of the template arguments
func TestContactHTML(t *testing.T) {
	var sb strings.Builder
	urls := template.ContactPageURLs{
		ContactForm: std_template.URL("/contact/form?Id=" + aContact.Id),
		ContactList: "/contact/list",
	}
	if err := template.WriteContactHTML(&sb, aContact, urls); err != nil {
		t.Fatal(err)
	}
	htmlDoc := sb.String()
	if _, err := html.Parse(strings.NewReader(htmlDoc)); err != nil {
		t.Errorf("invalid HTML: %v", err)
	}

	for name, value := range map[string]string{
		"Id":        aContact.Id.String(),
		"FirstName": aContact.FirstName,
		"LastName":  aContact.LastName,
		"Phone":     aContact.Phone,
		"Email":     aContact.Email,

		"ContactFormURL": string(urls.ContactForm),
		"ContactListURL": string(urls.ContactList),
	} {
		if !strings.Contains(htmlDoc, value) {
			t.Errorf("value %q of property %q not found in HTML", value, name)
		}
	}
}

// Check HTML is syntactically valid and that it contains all properties of the template arguments
func TestContactsHTML(t *testing.T) {
	var sb strings.Builder
	s := template.SearchPage{
		SearchTerm: "SEARCH_TERM",
		Contacts:   []contact.Contact{aContact},
	}
	if err := template.WriteContactsHTML(&sb, s); err != nil {
		t.Fatal(err)
	}
	htmlDoc := sb.String()
	if _, err := html.Parse(strings.NewReader(htmlDoc)); err != nil {
		t.Errorf("invalid HTML: %v", err)
	}

	for name, value := range map[string]string{
		"SearchTerm": s.SearchTerm,
		"Id":         aContact.Id.String(),
		"FirstName":  aContact.FirstName,
		"LastName":   aContact.LastName,
		"Phone":      aContact.Phone,
		"Email":      aContact.Email,
	} {
		if !strings.Contains(htmlDoc, value) {
			t.Errorf("value %q of property %q not found in HTML", value, name)
		}
	}
}

// Check HTML is syntactically valid and that it contains all properties of the template arguments
func TestContactFormHTML(t *testing.T) {
	var sb strings.Builder
	var f = template.ContactForm{
		Contact: aContact,
		Errors: templates.ErrorMap{
			"Email": fmt.Errorf("Invalid Email"),
		},
	}

	if err := template.WriteContactFormHTML(&sb, f); err != nil {
		t.Fatal(err)
	}
	htmlDoc := sb.String()
	if _, err := html.Parse(strings.NewReader(htmlDoc)); err != nil {
		t.Errorf("invalid HTML: %v", err)
	}

	for name, value := range map[string]string{
		"EmailErrorMessage": f.Errors["Email"].Error(),
		"Id":                aContact.Id.String(),
		"FirstName":         aContact.FirstName,
		"LastName":          aContact.LastName,
		"Phone":             aContact.Phone,
		"Email":             aContact.Email,
	} {
		if !strings.Contains(htmlDoc, value) {
			t.Errorf("value %q of property %q not found in HTML", value, name)
		}
	}
}
