//go:build integration

package test

func (s *Suite) Test_UpdateProfile() {
	id, err := s.profile.Create("John_Update", 25)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal(25, p.Age)

	err = s.profile.Update(id.String(), "John_Update", 26)
	s.NoError(err)

	p, err = s.profile.Get(id.String())
	s.NoError(err)

	s.Equal(26, p.Age)
}
