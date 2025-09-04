package httpclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func Example() {
	profile := New(Config{Host: "localhost", Port: "8080"})

	ctx := context.Background()

	id, err := create(ctx, profile)
	if err != nil {
		panic(fmt.Errorf("httpclient: example: create: %w", err))
	}

	p, err := profile.Get(ctx, id.String())
	if err != nil {
		panic(fmt.Errorf("httpclient: example: profile.Get_1: %w", err))
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

	err = update(ctx, profile, id.String())
	if err != nil {
		panic(fmt.Errorf("httpclient: example: update: %w", err))
	}

	p, err = profile.Get(ctx, id.String())
	if err != nil {
		panic(fmt.Errorf("httpclient: example: profile.Get_2: %w", err))
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

	err = profile.Delete(ctx, id.String())
	if err != nil {
		panic(fmt.Errorf("httpclient: example: profile.Delete: %w", err))
	}

	_, err = profile.Get(ctx, id.String())

	fmt.Println("The Example function completed successfully! Get request: ", err) //nolint:forbidigo
}

func create(ctx context.Context, profile *Client) (uuid.UUID, error) {
	requestCreate := CreateProfileRequest{
		Name:  "Andrey",
		Age:   37,
		Email: "7n1987@gmail.com",
		Phone: "+79634813074",
	}

	id, err := profile.Create(ctx, requestCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("profile.Create: %w", err)
	}

	return id, nil
}

func update(ctx context.Context, profile *Client, id string) error {
	var (
		name = "Anton"
		age  = 37
	)

	requestUpdate := UpdateProfileRequest{
		ID:   id,
		Name: &name,
		Age:  &age,
	}

	err := profile.Update(ctx, requestUpdate)
	if err != nil {
		return fmt.Errorf("profile.Update: %w", err)
	}

	return nil
}
