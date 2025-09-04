//go:build integration

package test

import (
	"context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"
)

func (s *Suite) Test_CreateProfile1() {
	request := httpclient.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(context.Background(), request)
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
	request := httpclient.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(context.Background(), request)
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
	request := httpclient.CreateProfileRequest{
		Name:  "",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err := s.profile.Create(context.Background(), request)
	s.ErrorContains(err, "validation")

	request = httpclient.CreateProfileRequest{
		Name:  "John_Create",
		Age:   17,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.profile.Create(context.Background(), request)
	s.ErrorContains(err, "validation")

	request = httpclient.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987gmail.com",
		Phone: "+79634813074",
	}

	_, err = s.profile.Create(context.Background(), request)
	s.ErrorContains(err, "validation")

	request = httpclient.CreateProfileRequest{
		Name:  "John_Create",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "79634813074",
	}

	_, err = s.profile.Create(context.Background(), request)
	s.ErrorContains(err, "validation")
}
