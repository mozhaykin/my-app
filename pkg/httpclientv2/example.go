package httpclientv2

import (
	"fmt"
)

func Example() { //nolint: funlen
	profile, err := New("http://localhost:8080/amozhaykin/my-app/api/v2")
	if err != nil {
		panic(err)
	}

	createRequest := CreateProfileRequest{
		Name:  "Andrey",
		Age:   37,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := profile.Create(createRequest)
	if err != nil {
		panic(err)
	}

	p, err := profile.Get(id.String())
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Profile:\n"+
			"  ID: 			%v\n"+
			"  CreatedAt: 	%v\n"+
			"  UpdatedAt: 	%v\n"+
			"  Name: 		%v\n"+
			"  Age: 		%v\n"+
			"  Status: 		%v\n"+
			"  Verifide:	%v\n"+
			"  Contacts: 	%v\n\n",
		p.ID, p.CreatedAt, p.UpdatedAt, p.Name, p.Age, p.Status, p.Verified, p.Contacts)

	var (
		name  = "Ekaterina"
		age   = 38
		email = "a7n1987@yandex.ru"
		phone = "+79634813069"
	)

	updateRequest := UpdateProfileRequest{
		ID:    id.String(),
		Name:  &name,
		Age:   &age,
		Email: &email,
		Phone: &phone,
	}

	err = profile.Update(updateRequest)
	if err != nil {
		panic(err)
	}

	p, err = profile.Get(id.String())
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Profile:\n"+
			"  ID: 			%v\n"+
			"  CreatedAt: 	%v\n"+
			"  UpdatedAt: 	%v\n"+
			"  Name: 		%v\n"+
			"  Age: 		%v\n"+
			"  Status: 		%v\n"+
			"  Verifide:	%v\n"+
			"  Contacts: 	%v\n\n",
		p.ID, p.CreatedAt, p.UpdatedAt, p.Name, p.Age, p.Status, p.Verified, p.Contacts)

	err = profile.Delete(id.String())
	if err != nil {
		panic(err)
	}

	_, err = profile.Get(id.String())

	fmt.Println("The Example function for httpclientv2 completed successfully! Get request:", err) //nolint:forbidigo
}
