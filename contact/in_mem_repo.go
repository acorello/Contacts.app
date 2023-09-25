package contact

type InMemoryRepository []Contact

func NewPopulatedInMemoryContactRepository() InMemoryRepository {
	return InMemoryRepository{
		{
			Id:        "0",
			FirstName: "Joe",
			LastName:  "Bloggs",
			Phone:     "+44(0)751123456",
			Email:     "joebloggs@example.com",
		},
	}
}

func (me InMemoryRepository) FindById(id string) (c Contact, found bool) {
	for _, c := range me {
		if c.Id == id {
			return c, true
		}
	}
	return c, false
}

func (me *InMemoryRepository) Delete(id string) {
	contacts := *me
	for i, c := range contacts {
		if c.Id == id {
			contacts[i] = contacts[len(contacts)-1]
			contacts[len(contacts)-1] = Contact{}
			*me = contacts[:len(contacts)-1]
			return
		}
	}
}

func (me InMemoryRepository) FindAll() (result []Contact) {
	for _, c := range me {
		result = append(result, c)
	}
	return
}

func (me *InMemoryRepository) Store(c Contact) error {
	for i, x := range *me {
		if x.Id == c.Id {
			(*me)[i] = c
			return nil
		}
	}
	*me = append(*me, c)
	return nil
}

func (me InMemoryRepository) FindBySearchTerm(term string) (result []Contact) {
	for _, c := range me {
		if c.AnyFieldContains(term) {
			result = append(result, c)
		}
	}
	return
}
