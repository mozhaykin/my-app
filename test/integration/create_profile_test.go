//go:build integration

package test

func (s *Suite) Test_CreateProfile() {
	id, err := s.profile.Create("John_Create", 25)
	s.NoError(err)

	p, err := s.profile.Get(id.String())
	s.NoError(err)

	s.Equal("John_Create", p.Name)
	s.Equal(25, p.Age)
}
