//go:build integration

package test

func (s *Suite) Test_CreateProfile() {
	request := CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.client.Create(ctx, request)
	s.NoError(err)

	p, err := s.client.Get(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Create", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_CreateProfile_IsInvalid() {
	request := CreateProfileRequest{
		Name:  "",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err := s.client.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = CreateProfileRequest{
		Name:  "John_Create",
		Age:   17,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.client.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.client.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "invalid-phone",
	}

	_, err = s.client.Create(ctx, request)
	s.ErrorContains(err, "validation")
}
