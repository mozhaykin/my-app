//go:build integration

package test

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv1"
)

func (s *Suite) Test_CreateProfile() {
	request := httpclientv1.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(ctx, request)
	s.NoError(err)

	p, err := s.profile.Get(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Create", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_CreateProfile_IsInvalid() {
	request := httpclientv1.CreateProfileRequest{
		Name:  "",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err := s.profile.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = httpclientv1.CreateProfileRequest{
		Name:  "John_Create",
		Age:   17,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.profile.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = httpclientv1.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.profile.Create(ctx, request)
	s.ErrorContains(err, "validation")

	request = httpclientv1.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "79634813074",
	}

	_, err = s.profile.Create(ctx, request)
	s.ErrorContains(err, "validation")
}
