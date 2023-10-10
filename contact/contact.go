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

type Page struct {
	Offset, Size int
}

func (me Page) Next() Page {
	return Page{
		Offset: me.Offset + 1,
		Size:   me.Size,
	}
}

func (me Page) StartOffset() int {
	return me.Offset * me.Size
}

func (me Page) EndOffset() int {
	return me.StartOffset() + me.Size
}

type Repository interface {
	FindById(id Id) (c Contact, found bool)
	Delete(id Id)
	FindAll(page Page) (result []Contact, more bool)
	Store(c Contact) error
	FindBySearchTerm(term string, page Page) (result []Contact, more bool)
	FindIdByEmail(email string) (res Id, found bool)
}
