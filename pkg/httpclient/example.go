package httpclient

import (
	"context"
	"fmt"
)

func Example() {
	profile := New(Config{Host: "k8s.goscl.ru", Port: "8080"})

	ctx := context.Background()

	id, err := profile.Create(ctx, "Andrey", 37, "7n1987@gmail.com", "+79634813074")
	if err != nil {
		panic(fmt.Errorf("profile.Create: %w", err))
	}

	p, err := profile.Get(ctx, id.String())
	if err != nil {
		panic(fmt.Errorf("profile.Get1: %w", err))
	}

	fmt.Printf( //nolint forbidigo
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

	err = profile.Update(ctx, id.String(), "Anton", 30, "", "")
	if err != nil {
		panic(fmt.Errorf("profile.Update: %w", err))
	}

	p, err = profile.Get(ctx, id.String())
	if err != nil {
		panic(fmt.Errorf("profile.Get2: %w", err))
	}

	fmt.Printf( //nolint forbidigo
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
		panic(fmt.Errorf("profile.Delete: %w", err))
	}

	_, err = profile.Get(ctx, id.String())

	fmt.Println("The Example function completed successfully! Get request: ", err) //nolint:forbidigo
}
