package main

import (
	"fmt"
	"testing"
)

type Customer struct{}

var count uint32

func (c *Customer) OnDirty(i interface{}) {
	count++
	fmt.Println("OnDirty", i, count)
}

func TestExample(t *testing.T) {
	observer := &Customer{}

	user := NewUser()
	user.Attach(observer)
	user.SetScore(18) // dirty 1

	p := NewPerson()
	p.SetName("123") // no dirty

	user.SetBaseInfo(p) // dirty 2

	p.SetAge(3) // dirty 3

	persons := make([]*Person, 0)
	p1 := NewPerson()
	p1.SetName("p1")
	persons = append(persons, p1)

	p2 := NewPerson()
	p2.SetName("p2")
	persons = append(persons, p2)

	wpersons := NewWrapPersonFriendsFromSlice(persons)
	p.SetFriends(wpersons) // dirty 4

	p1.SetAge(1) // dirty 5
	p2.SetAge(2) // dirty 6

	p3 := NewPerson()
	p3.SetName("p3")
	wpersons.Append(p3) // dirty 7

	p3.SetAge(3) // dirty 8

	p4 := NewPerson()
	p4.SetName("p4")
	persons = append(persons, p2) // no dirty
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。

	peoples := make(map[string]*Person)
	pp1 := NewPerson()
	peoples["pp1"] = pp1

	wpeoples := NewWrapPersonPeoplesFromMap(peoples)
	p.SetPeoples(wpeoples) // dirty 9
	pp1.SetAge(11)         // dirty 10

	pp2 := NewPerson()
	wpeoples.Set("pp2", pp2) // dirty 11
	pp2.SetAge(22)           // dirty 12

	pp3 := NewPerson()
	peoples["pp3"] = pp3 // no dirty
	pp3.SetAge(33)       // no dirty
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。
}
