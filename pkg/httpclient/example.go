package httpclient

import "fmt"

func Example() {
	profile := New("k8s.goscl.ru")

	id, err := profile.Create("John", 25)
	if err != nil {
		panic(err)
	}

	p, err := profile.Get(id.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(p.ID)
	fmt.Println(p.Name)
	fmt.Println(p.Age)

	err = profile.Update(id.String(), "John Doe", 26)
	if err != nil {
		panic(err)
	}

	p, err = profile.Get(id.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(p.ID)
	fmt.Println(p.Name)
	fmt.Println(p.Age)

	err = profile.Delete(id.String())
	if err != nil {
		panic(err)
	}

	_, err = profile.Get(id.String())

	fmt.Println("Get request: ", err)
}
