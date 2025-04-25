package domain

import "github.com/google/uuid"

type Name string

type Age int

type Profile struct {
	ID   uuid.UUID `json:"id"`
	Name Name      `json:"name"`
	Age  Age       `json:"age"`
}

func NewProfile(name string, age int, id ...uuid.UUID) (Profile, error) {
	var p Profile

	if name == "" {
		return p, ErrEmptyName
	}

	if age < 18 {
		return p, ErrAgeLessThan18
	}

	p = Profile{
		Name: Name(name),
		Age:  Age(age),
	}

	if len(id) > 0 {
		p.ID = id[0]
	} else {
		p.ID = uuid.New()
	}

	return p, nil
}
