package domain

type Name string

type Age int

type Profile struct {
	Name Name
	Age  Age
}

func NewProfile(name string, age int) (Profile, error) {
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

	return p, nil
}
