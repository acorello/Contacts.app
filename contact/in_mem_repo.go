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
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000002"),
				FirstName: "Jane",
				LastName:  "Doe",
				Phone:     "+44(0)751123457",
				Email:     "janedoe@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000003"),
				FirstName: "Sam",
				LastName:  "Smith",
				Phone:     "+44(0)751123458",
				Email:     "samsmith@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000004"),
				FirstName: "Ann",
				LastName:  "Taylor",
				Phone:     "+44(0)751123459",
				Email:     "anntaylor@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000005"),
				FirstName: "Bob",
				LastName:  "Brown",
				Phone:     "+44(0)751123460",
				Email:     "bobbrown@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000006"),
				FirstName: "Lucy",
				LastName:  "Green",
				Phone:     "+44(0)751123461",
				Email:     "lucygreen@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000007"),
				FirstName: "Dan",
				LastName:  "White",
				Phone:     "+44(0)751123462",
				Email:     "danwhite@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000008"),
				FirstName: "Eva",
				LastName:  "Black",
				Phone:     "+44(0)751123463",
				Email:     "evablack@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000009"),
				FirstName: "Tom",
				LastName:  "Gray",
				Phone:     "+44(0)751123464",
				Email:     "tomgray@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000010"),
				FirstName: "Sue",
				LastName:  "Jones",
				Phone:     "+44(0)751123465",
				Email:     "suejones@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000011"),
				FirstName: "Lee",
				LastName:  "Davis",
				Phone:     "+44(0)751123466",
				Email:     "leedavis@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000012"),
				FirstName: "Amy",
				LastName:  "Adams",
				Phone:     "+44(0)751123467",
				Email:     "amyadams@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000013"),
				FirstName: "Max",
				LastName:  "Mills",
				Phone:     "+44(0)751123468",
				Email:     "maxmills@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000014"),
				FirstName: "Tina",
				LastName:  "Turner",
				Phone:     "+44(0)751123469",
				Email:     "tinaturner@example.com",
			},
			{
				Id:        MustParseId("00000000-0000-0000-0000-000000000015"),
				FirstName: "Rob",
				LastName:  "Rider",
				Phone:     "+44(0)751123470",
				Email:     "robrider@example.com",
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
