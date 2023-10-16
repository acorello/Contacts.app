package contact

import (
	"fmt"
	"log"
	"slices"
)

type InMemoryRepository struct {
	contacts []Contact
}

func NewPopulatedInMemoryContactRepository() InMemoryRepository {
	return InMemoryRepository{
		contacts: slices.Clone(fixedContactsList),
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

func (me InMemoryRepository) FindAll(page Page) (result []Contact, more bool) {
	maxEnd := len(me.contacts)
	if page.StartOffset() > maxEnd {
		return nil, false
	}
	start := page.StartOffset()
	pageEnd := page.EndOffset()
	end := min(pageEnd, maxEnd)
	result = slices.Clone(me.contacts[start:end])
	return result, maxEnd > pageEnd
}

// TODO: implement validation ( eg. [e-mail]--N--1--[contactId] )
func (me *InMemoryRepository) Store(c Contact) error {
	log.Printf("Storing %#v", c)
	if err := me.checkEmailOwner(c); err != nil {
		return err
	}
	existingIdx := slices.IndexFunc(me.contacts, c.Id.HasSameId)
	if existingIdx >= 0 {
		me.contacts[existingIdx] = c
	} else {
		me.contacts = append(me.contacts, c)
	}
	return nil
}

func (me *InMemoryRepository) checkEmailOwner(c Contact) error {
	var alreadyAssignedId, found = me.FindIdByEmail(c.Email)
	if found && c.Id != alreadyAssignedId {
		return fmt.Errorf("e-mail already assigned to contact with id %q", alreadyAssignedId)
	}
	return nil
}

func (me InMemoryRepository) FindBySearchTerm(term string, page Page) (result []Contact, more bool) {
	// me.contacts.findBy(p).drop(page.StartOffset()).take(page.Size)
	start := page.StartOffset()
	foundCount := 0
	size := page.Size + 1 // we try fetching one more to tell if there is anothe page
	for _, c := range me.contacts {
		if len(result) >= size {
			break
		}
		if c.AnyFieldContains(term) {
			if foundCount >= start {
				result = append(result, c)
			}
			foundCount += 1
		}
	}
	if len(result) == size {
		return result[:len(result)-1], true
	} else {
		return result, false
	}
}
