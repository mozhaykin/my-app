//go:build integration

package test

func (s *Suite) Test_DeleteProfile() {
	request := CreateProfileRequest{
		Name:  "John_Delete",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.client.Create(ctx, request)
	s.NoError(err)

	p, err := s.client.Get(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Delete", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)

	err = s.client.Delete(ctx, id.String())
	s.NoError(err)

	_, err = s.client.Get(ctx, id.String())
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_DeleteProfile_NotFound() {
	err := s.client.Delete(ctx, "e6799c89-c560-45a2-a3da-b3f1eb9bee2b")
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_DeleteProfile_UuidIsInvalid() {
	err := s.client.Delete(ctx, "c6799c89c560-45a2-a3da-b3f1eb9bee2b")
	s.ErrorContains(err, "uuid is invalid")
}
