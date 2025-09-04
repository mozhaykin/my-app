//go:build integration

package test

import (
	"context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"
)

func (s *Suite) Test_UpdateProfile() {
	requestCreate := httpclient.CreateProfileRequest{
		Name:  "John_Update",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(context.Background(), requestCreate)
	s.NoError(err)

	p, err := s.profile.Get(context.Background(), id.String())
	s.NoError(err)

	s.Equal("John_Update", string(p.Name))
	s.Equal(25, int(p.Age))
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, int(p.Status))
	s.Equal(false, p.Verified)

	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclient.UpdateProfileRequest{
		ID:    id.String(),
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err = s.profile.Update(context.Background(), requestUpdate)
	s.NoError(err)

	p, err = s.profile.Get(context.Background(), id.String())
	s.NoError(err)

	s.Equal("New_John_Update", string(p.Name))
	s.Equal(26, int(p.Age))
	s.Equal("a7n1987@yandex.ru", p.Contacts.Email)
	s.Equal("+79634813069", p.Contacts.Phone)
	s.Equal(1, int(p.Status))
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_UpdateProfile_NotFound() {
	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclient.UpdateProfileRequest{
		ID:    "c6799c89-c560-45a2-a3da-b3f1eb9bee2b",
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err := s.profile.Update(context.Background(), requestUpdate)
	s.EqualError(err, httpclient.ErrNotFound.Error())
}

func (s *Suite) Test_UpdateProfile_UUIDInvalid() {
	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclient.UpdateProfileRequest{
		ID:    "c6799c89c560-45a2-a3da-b3f1eb9bee2b",
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err := s.profile.Update(context.Background(), requestUpdate)
	s.ErrorContains(err, "uuid is invalid")
}

func (s *Suite) Test_UpdateProfile_NoChangesFound() {
	requestCreate := httpclient.CreateProfileRequest{
		Name:  "Ben_Update",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(context.Background(), requestCreate)
	s.NoError(err)

	var (
		newName  = "Ben_Update"
		newAge   = 25
		newEmail = "7n1987@gmail.com"
		newPhone = "+79634813074"
	)

	requestUpdate := httpclient.UpdateProfileRequest{
		ID:    id.String(),
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err = s.profile.Update(context.Background(), requestUpdate)
	s.ErrorContains(err, "no changes found")
}

func (s *Suite) Test_UpdateProfile_AllFieldsAreEmpty() {
	requestUpdate := httpclient.UpdateProfileRequest{
		ID:    "c6799c89-c560-45a2-a3da-b3f1eb9bee2b",
		Name:  nil,
		Age:   nil,
		Email: nil,
		Phone: nil,
	}

	err := s.profile.Update(context.Background(), requestUpdate)
	s.ErrorContains(err, "all fields for update are empty")
}
