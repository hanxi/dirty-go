package dirty_tmpl

type Person struct {
	name    string
	age     int
	friends []*Person
	peoples map[string]*Person
}

type Man struct {
	baseInfo *Person
	score    uint32
}
