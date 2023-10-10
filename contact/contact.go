package contact

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Id string

func NewId() Id {
	return Id(uuid.NewString())
}

func ParseId(s string) (Id, error) {
	u, err := uuid.Parse(s)
	return Id(u.String()), err
}

func MustParseId(s string) Id {
	u, err := uuid.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseId failed for %q: %v", s, err))
	}
	return Id(u.String())
}

func (me Id) String() string {
	return string(me)
}

func (me Id) HasSameId(c Contact) bool {
	return me == c.Id
}

type Contact struct {
	Id
	FirstName, LastName, Phone, Email string
}

func (my Contact) AnyFieldContains(s string) bool {
	p := strings.Contains
	return p(my.FirstName, s) || p(my.LastName, s) || p(my.Phone, s) || p(my.Email, s)
}

type Repository interface {
	FindById(id Id) (c Contact, found bool)
	Delete(id Id)
	FindAll() (result []Contact)
	Store(c Contact) error
	FindBySearchTerm(term string) (result []Contact)
	FindIdByEmail(email string) (res Id, found bool)
}
