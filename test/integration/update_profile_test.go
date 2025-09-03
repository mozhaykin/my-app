//go:build integration

package test

import (
	"context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"
)

func (s *Suite) Test_UpdateProfile() {
	id, err := s.profile.Create(context.Background(), "John_Update", 25, "7n1987@gmail.com", "+79634813074")
	s.NoError(err)

	p, err := s.profile.Get(context.Background(), id.String())
	s.NoError(err)

	s.Equal("John_Update", string(p.Name))
	s.Equal(25, int(p.Age))
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)
	s.Equal(1, int(p.Status))
	s.Equal(false, p.Verified)

	err = s.profile.Update(context.Background(), id.String(), "New_John_Update", 26, "a7n1987@yandex.ru", "+79634813069")
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

func (s *Suite) Test_UpdateProfile_ErrUUIDIsEmpty() {
	err := s.profile.Update(context.Background(), "", "New_John_Update", 26, "a7n1987@yandex.ru", "+79634813069")
	s.ErrorContains(err, "uuid is empty")
}

func (s *Suite) Test_UpdateProfile_NotFound() {
	err := s.profile.Update(context.Background(), "c6799c89-c560-45a2-a3da-b3f1eb9bee2b", "New_John_Update", 26, "a7n1987@yandex.ru", "+79634813069")
	s.EqualError(err, httpclient.ErrNotFound.Error())
}

func (s *Suite) Test_UpdateProfile_UUIDInvalid() {
	err := s.profile.Update(context.Background(), "c6799c89c560-45a2-a3da-b3f1eb9bee2b", "New_John_Update", 26, "a7n1987@yandex.ru", "+79634813069")
	s.ErrorContains(err, "uuid is invalid")
}

func (s *Suite) Test_UpdateProfile_NoChangesFound() {
	id, err := s.profile.Create(context.Background(), "Ben_Update", 25, "7n1987@gmail.com", "+79634813074")
	s.NoError(err)

	err = s.profile.Update(context.Background(), id.String(), "Ben_Update", 25, "7n1987@gmail.com", "+79634813074")
	s.ErrorContains(err, "no changes found")
}
