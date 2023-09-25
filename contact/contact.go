package contact

import (
	"strings"
)

type Contact struct {
	Id, FirstName, LastName, Phone, Email string
}

func (my Contact) AnyFieldContains(s string) bool {
	p := strings.Contains
	return p(my.FirstName, s) || p(my.LastName, s) || p(my.Phone, s) || p(my.Email, s)
}

type Repository interface {
	FindById(id string) (c Contact, found bool)
	Delete(id string)
	FindAll() (result []Contact)
	Store(c Contact) error
	FindBySearchTerm(term string) (result []Contact)
}
