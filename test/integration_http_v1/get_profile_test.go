//go:build integration

package test

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv1"
)

func (s *Suite) Test_GetProfile_Ok() {
	request := httpclientv1.CreateProfileRequest{
		Name:  "John_Get",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(ctx, request)
	s.NoError(err)

	p, err := s.profile.Get(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Get", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_GetProfile_NotFound() {
	_, err := s.profile.Get(ctx, "c6799c89-c560-45a2-afda-b3f1eb9bee2b")
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_GetProfile_UuidIsInvalid() {
	_, err := s.profile.Get(ctx, "c6799c89c560-45a2-a3da-b3f1eb9bee2b")
	s.ErrorContains(err, "uuid is invalid")
}
