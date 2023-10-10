package contact

import (
	"log"
	"slices"
)

type InMemoryRepository struct {
	contacts []Contact
}

func NewPopulatedInMemoryContactRepository() InMemoryRepository {
	return InMemoryRepository{
		contacts: []Contact{
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000001"),
				FirstName: "Joe",
				LastName:  "Bloggs",
				Phone:     "+44(0)751123456",
				Email:     "joebloggs@example.com",
			},
		},
	}
}

func (me InMemoryRepository) FindById(id Id) (c Contact, found bool) {
	idx := slices.IndexFunc(me.contacts, id.HasSameId)
	if idx >= 0 {
		return me.contacts[idx], true
	} else {
		return c, false
	}
}

func (me InMemoryRepository) FindIdByEmail(email string) (res Id, found bool) {
	for i := range me.contacts {
		if me.contacts[i].Email == email {
			return me.contacts[i].Id, true
		}
	}
	var zeroId Id
	return zeroId, false
}

func (me *InMemoryRepository) Delete(id Id) {
	me.contacts = slices.DeleteFunc(me.contacts, id.HasSameId)
}

func (me InMemoryRepository) FindAll() []Contact {
	return slices.Clone(me.contacts)
}

func (me *InMemoryRepository) Store(c Contact) error {
	log.Printf("Storing %#v", c)
	existingIdx := slices.IndexFunc(me.contacts, c.Id.HasSameId)
	if existingIdx >= 0 {
		me.contacts[existingIdx] = c
	} else {
		me.contacts = append(me.contacts, c)
	}
	return nil
}

func (me InMemoryRepository) FindBySearchTerm(term string) (result []Contact) {
	for _, c := range me.contacts {
		if c.AnyFieldContains(term) {
			result = append(result, c)
		}
	}
	return
}
