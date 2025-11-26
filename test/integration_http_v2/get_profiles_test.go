//go:build integration

package test

func (s *Suite) Test_GetProfiles_Ok() {
	requestCreate1 := CreateProfileRequest{
		Name:  "John1_Get",
		Age:   25,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	requestCreate2 := CreateProfileRequest{
		Name:  "John2_Get",
		Age:   24,
		Email: "1987@gmail.com",
		Phone: "+79634813069",
	}

	_, err := s.client.Create(ctx, requestCreate1)
	s.NoError(err)

	_, err = s.client.Create(ctx, requestCreate2)
	s.NoError(err)

	requestGet := GetProfilesRequest{
		Sort:   "name",
		Order:  "asc",
		Offset: 0,
		Limit:  10,
	}

	profiles, err := s.client.GetProfiles(ctx, requestGet)
	s.NoError(err)

	s.Equal(2, len(profiles))

	p := profiles[0]

	s.Equal("John1_Get", p.Name)
	s.Equal(25, p.Age)
	s.Equal("7n1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813074", p.Contacts.Phone)

	p = profiles[1]

	s.Equal("John2_Get", p.Name)
	s.Equal(24, p.Age)
	s.Equal("1987@gmail.com", p.Contacts.Email)
	s.Equal("+79634813069", p.Contacts.Phone)
}

func (s *Suite) Test_GetProfiles_NotFound() {
	requestGet := GetProfilesRequest{
		Sort:   "name",
		Order:  "asc",
		Offset: 0,
		Limit:  10,
	}
	_, err := s.client.GetProfiles(ctx, requestGet)
	s.ErrorContains(err, "not found")
}

func (s *Suite) Test_GetProfiles_Input_IsInvalid() {
	requestGet := GetProfilesRequest{
		Sort:   "invalid",
		Order:  "asc",
		Offset: 0,
		Limit:  10,
	}
	_, err := s.client.GetProfiles(ctx, requestGet)
	s.ErrorContains(err, "validation")
}
