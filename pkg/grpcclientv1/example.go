package grpcclientv1

import (
	"fmt"
)

func Example() { //nolint: funlen
	profile, err := New("localhost:50051")
	if err != nil {
		panic(err)
	}

	defer profile.Close()

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
			"  DeletedAt: 	%v\n"+
			"  Name: 		%v\n"+
			"  Age: 		%v\n"+
			"  Status: 		%v\n"+
			"  Verifide:	%v\n"+
			"  Contacts: 	%v\n\n",
		p.ID, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.Name, p.Age, p.Status, p.Verified, p.Contacts)

	var (
		name = "Anton"
		age  = 37
	)

	updateRequest := UpdateProfileRequest{
		ID:   id.String(),
		Name: &name,
		Age:  &age,
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
			"  DeletedAt: 	%v\n"+
			"  Name: 		%v\n"+
			"  Age: 		%v\n"+
			"  Status: 		%v\n"+
			"  Verifide:	%v\n"+
			"  Contacts: 	%v\n\n",
		p.ID, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.Name, p.Age, p.Status, p.Verified, p.Contacts)

	err = profile.Delete(id.String())
	if err != nil {
		panic(err)
	}

	_, err = profile.Get(id.String())

	fmt.Println("The Example function for grpcclientv1 completed successfully! Get request:", err) //nolint:forbidigo
}
