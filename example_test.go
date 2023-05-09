package main

import (
	"fmt"
	"testing"

	"github.com/hanxi/dirty-go/dirty_out"
)

type Customer struct{}

var count uint32

func (c *Customer) OnDirty(i interface{}) {
	count++
	fmt.Println("OnDirty", i, count)
}

func TestExample(t *testing.T) {
	observer := &Customer{}

	man := dirty_out.NewMan()
	man.Attach(observer)
	man.SetScore(18) // dirty 1

	p := dirty_out.NewPerson()
	p.SetName("123") // no dirty

	man.SetBaseInfo(p) // dirty 2

	p.SetAge(3) // dirty 3

	persons := make([]*dirty_out.Person, 0)
	p1 := dirty_out.NewPerson()
	p1.SetName("p1")
	persons = append(persons, p1)

	p2 := dirty_out.NewPerson()
	p2.SetName("p2")
	persons = append(persons, p2)

	wpersons := dirty_out.NewArrPersonFriendsFromSlice(persons)
	p.SetFriends(wpersons) // dirty 4

	p1.SetAge(1) // dirty 5
	p2.SetAge(2) // dirty 6

	p3 := dirty_out.NewPerson()
	p3.SetName("p3")
	wpersons.Append(p3) // dirty 7

	p3.SetAge(3) // dirty 8

	p4 := dirty_out.NewPerson()
	p4.SetName("p4")
	persons = append(persons, p2) // no dirty
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。

	peoples := make(map[string]*dirty_out.Person)
	pp1 := dirty_out.NewPerson()
	peoples["pp1"] = pp1

	wpeoples := dirty_out.NewMapPersonPeoplesFromMap(peoples)
	p.SetPeoples(wpeoples) // dirty 9
	pp1.SetAge(11)         // dirty 10

	pp2 := dirty_out.NewPerson()
	wpeoples.Set("pp2", pp2) // dirty 11
	pp2.SetAge(22)           // dirty 12

	pp3 := dirty_out.NewPerson()
	peoples["pp3"] = pp3 // no dirty
	pp3.SetAge(33)       // no dirty
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。
}
