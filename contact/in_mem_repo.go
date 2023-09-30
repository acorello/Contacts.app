package contact

import "slices"

type InMemoryRepository struct {
	contacts []Contact
}

func NewPopulatedInMemoryContactRepository() InMemoryRepository {
	return InMemoryRepository{
		contacts: []Contact{
			{
				Id:        "0",
				FirstName: "Joe",
				LastName:  "Bloggs",
				Phone:     "+44(0)751123456",
				Email:     "joebloggs@example.com",
			},
		},
	}
}

func (me InMemoryRepository) FindById(id Id) (c Contact, found bool) {
	idx := slices.IndexFunc(me.contacts, id.Has)
	if idx >= 0 {
		return me.contacts[idx], true
	} else {
		return c, false
	}
}

func (me InMemoryRepository) Delete(id Id) {
	me.contacts = slices.DeleteFunc(me.contacts, id.Has)
}

func (me InMemoryRepository) FindAll() (result []Contact) {
	for _, c := range me.contacts {
		result = append(result, c)
	}
	return
}

func (me InMemoryRepository) Store(c Contact) error {
	existingIdx := slices.IndexFunc(me.contacts, c.Id.Has)
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
