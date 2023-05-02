package main

type Person struct {
	Base
	name    string
	age     int
	friends []*Person
	peoples map[string]*Person
}

type User struct {
	Base
	baseInfo *Person
	score    uint32
}
