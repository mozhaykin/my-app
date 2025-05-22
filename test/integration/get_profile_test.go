//go:build integration

package test

import "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"

func (s *Suite) Test_GetProfile_Ok() {
	id, err := s.profile.Create("John_Get", 25)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("John_Get", p.Name)
	s.Equal(25, p.Age)
}

func (s *Suite) Test_GetProfile_NotFound() {
	_, err := s.profile.Get("c6799c89-c560-45a2-afda-b3f1eb9bee2b")
	s.EqualError(err, httpclient.ErrNotFound.Error())
}
