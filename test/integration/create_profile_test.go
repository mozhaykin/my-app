//go:build integration

package test

import "context"

func (s *Suite) Test_CreateProfile1() {
	id, err := s.profile.Create(context.Background(), "John_Create", 25, "7n1987@gmail.com", "+79634813074")
	s.NoError(err)

	p, err := s.profile.Get(context.Background(), id.String())
	s.NoError(err)

	s.Equal("John_Create", string(p.Name))
	s.Equal(25, int(p.Age))
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, int(p.Status))
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_CreateProfile2() {
	id, err := s.profile.Create(context.Background(), "John_Create", 25, "7n1987@gmail.com", "+79634813074")
	s.NoError(err)

	p, err := s.profile.Get(context.Background(), id.String())
	s.NoError(err)

	s.Equal("John_Create", string(p.Name))
	s.Equal(25, int(p.Age))
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, int(p.Status))
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_CreateProfile_IsInvalid() {
	_, err := s.profile.Create(context.Background(), "", 25, "7n1987@gmail.com", "+79634813074")
	s.ErrorContains(err, "validation")

	_, err = s.profile.Create(context.Background(), "John_Create", 17, "7n1987@gmail.com", "+79634813074")
	s.ErrorContains(err, "validation")

	_, err = s.profile.Create(context.Background(), "John_Create", 25, "7n1987gmail.com", "+79634813074")
	s.ErrorContains(err, "validation")

	_, err = s.profile.Create(context.Background(), "John_Create", 25, "7n1987@gmail.com", "74813074")
	s.ErrorContains(err, "validation")
}
