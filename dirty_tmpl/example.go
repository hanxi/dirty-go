package dirty_tmpl

type Person struct {
	Name    string
	Age     int
	Friends []*Person
	Peoples map[string]*Person
}

type Man struct {
	BaseInfo *Person
	Score    uint32
}
