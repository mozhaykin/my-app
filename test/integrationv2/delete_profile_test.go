//go:build integration

package test

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv2"
)

func (s *Suite) Test_DeleteProfile() {
	request := httpclientv2.CreateProfileRequest{
		Name:  "John_Delete",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}
	id, err := s.profile.Create(request)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("John_Delete", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)

	err = s.profile.Delete(id.String())
	s.NoError(err)

	_, err = s.profile.Get(id.String())
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_DeleteProfile_NotFound() {
	err := s.profile.Delete("e6799c89-c560-45a2-a3da-b3f1eb9bee2b")
	s.ErrorContains(err, "not found")
}
