//go:build integration

package test

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv1"
)

func (s *Suite) Test_UpdateProfile() {
	requestCreate := httpclientv1.CreateProfileRequest{
		Name:  "John_Update",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(requestCreate)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("John_Update", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)

	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclientv1.UpdateProfileRequest{
		ID:    id.String(),
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err = s.profile.Update(requestUpdate)
	s.NoError(err)

	p, err = s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("New_John_Update", p.Name)
	s.Equal(26, p.Age)
	s.Equal("a7n1987@yandex.ru", p.Contacts.Email)
	s.Equal("+79634813069", p.Contacts.Phone)
	s.Equal(1, p.Status)
	s.Equal(false, p.Verified)
}

func (s *Suite) Test_UpdateProfile_NotFound() {
	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclientv1.UpdateProfileRequest{
		ID:    "c6799c89-c560-45a2-a3da-b3f1eb9bee2b",
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err := s.profile.Update(requestUpdate)
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_UpdateProfile_NoChangesFound() {
	requestCreate := httpclientv1.CreateProfileRequest{
		Name:  "Ben_Update",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := s.profile.Create(requestCreate)
	s.NoError(err)

	var (
		newName  = "Ben_Update"
		newAge   = 25
		newEmail = "7n1987@gmail.com"
		newPhone = "+79634813074"
	)

	requestUpdate := httpclientv1.UpdateProfileRequest{
		ID:    id.String(),
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err = s.profile.Update(requestUpdate)
	s.ErrorContains(err, "no changes found")
}

func (s *Suite) Test_UpdateProfile_AllFieldsAreEmpty() {
	requestUpdate := httpclientv1.UpdateProfileRequest{
		ID:    "c6799c89-c560-45a2-a3da-b3f1eb9bee2b",
		Name:  nil,
		Age:   nil,
		Email: nil,
		Phone: nil,
	}

	err := s.profile.Update(requestUpdate)
	s.ErrorContains(err, "all fields for update are empty")
}

func (s *Suite) Test_UpdateProfile_UUIDInvalid() {
	var (
		newName  = "New_John_Update"
		newAge   = 26
		newEmail = "a7n1987@yandex.ru"
		newPhone = "+79634813069"
	)

	requestUpdate := httpclientv1.UpdateProfileRequest{
		ID:    "c6799c89c560-45a2-a3da-b3f1eb9bee2b",
		Name:  &newName,
		Age:   &newAge,
		Email: &newEmail,
		Phone: &newPhone,
	}

	err := s.profile.Update(requestUpdate)
	s.ErrorContains(err, "uuid is invalid")
}
