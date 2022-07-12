package test

type Person struct {
	ID   int    `json:"id" seal:"id"`
	Name string `json:"name" seal:"name"`
	Age  int    `json:"age" seal:"age"`
}

func (p *Person) Zero() bool {
	return p.ID == 0
}
