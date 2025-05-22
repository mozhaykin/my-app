//go:build integration

package test

import "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"

func (s *Suite) Test_DeleteProfile() {
	id, err := s.profile.Create("John_Delete", 25)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("John_Delete", p.Name)
	s.Equal(25, p.Age)

	err = s.profile.Delete(id.String())
	s.NoError(err)

	_, err = s.profile.Get(id.String())
	s.EqualError(err, httpclient.ErrNotFound.Error())
}
