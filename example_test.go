package main

import (
	"encoding/json"
	"runtime/debug"
	"testing"

	"github.com/hanxi/dirty-go/dirty_out"
)

type Customer struct{}

var count uint32

func (c *Customer) OnDirty(i interface{}) {
	count++
}

func ResetCount() {
	count = 0
}

func CheckCount(t *testing.T, want uint32) {
	if count != want {
		t.Errorf("Expected count to be %d, but got %d, stack:%s", want, count, debug.Stack())
	}
}

func TestExample(t *testing.T) {
	observer := &Customer{}

	ResetCount()
	man := dirty_out.NewMan()
	man.Attach(observer)
	man.SetScore(18) // dirty
	CheckCount(t, 1)

	ResetCount()
	p := dirty_out.NewPerson()
	p.SetName("123") // no dirty
	CheckCount(t, 0)

	ResetCount()
	man.SetBaseInfo(p) // dirty
	CheckCount(t, 1)

	ResetCount()
	p.SetAge(3) // dirty
	CheckCount(t, 1)

	ResetCount()
	persons := make([]*dirty_out.Person, 0)
	p1 := dirty_out.NewPerson()
	p1.SetName("p1")
	persons = append(persons, p1)
	p2 := dirty_out.NewPerson()
	p2.SetName("p2")
	persons = append(persons, p2)
	CheckCount(t, 0)

	ResetCount()
	wpersons := dirty_out.NewArrPersonFriendsFromSlice(persons)
	p.SetFriends(wpersons) // dirty
	CheckCount(t, 1)

	ResetCount()
	p1.SetAge(1) // dirty
	p2.SetAge(2) // dirty
	CheckCount(t, 2)

	ResetCount()
	p3 := dirty_out.NewPerson()
	p3.SetName("p3")
	wpersons.Append(p3) // dirty
	CheckCount(t, 1)

	ResetCount()
	p3.SetAge(3) // dirty
	CheckCount(t, 1)

	ResetCount()
	p4 := dirty_out.NewPerson()
	p4.SetName("p4")
	persons = append(persons, p2) // no dirty
	CheckCount(t, 0)
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。

	ResetCount()
	peoples := make(map[string]*dirty_out.Person)
	pp1 := dirty_out.NewPerson()
	peoples["pp1"] = pp1
	CheckCount(t, 0)

	ResetCount()
	wpeoples := dirty_out.NewMapPersonPeoplesFromMap(peoples)
	p.SetPeoples(wpeoples) // dirty
	CheckCount(t, 1)
	pp1.SetAge(11) // dirty
	CheckCount(t, 2)

	ResetCount()
	pp2 := dirty_out.NewPerson()
	wpeoples.Set("pp2", pp2) // dirty
	CheckCount(t, 1)
	pp2.SetAge(22) // dirty
	CheckCount(t, 2)

	ResetCount()
	pp3 := dirty_out.NewPerson()
	peoples["pp3"] = pp3 // no dirty
	pp3.SetAge(33)       // no dirty
	CheckCount(t, 0)
	// 这种情况没监听到修改影响不是很大，因为修改的临时数据，不是存储数据，
	// 不需要重启服务器就可以很容易的测试出问题，表现就是操作后数据没变化。
}

func TestNotifyDirty(t *testing.T) {
	observer := &Customer{}

	user := dirty_out.NewUser()
	user.Attach(observer)

	// 测试 SetName 方法是否触发了 NotifyDirty
	ResetCount()
	user.SetName("Alice")
	CheckCount(t, 1)

	// 测试 SetAge 方法是否触发了 NotifyDirty
	ResetCount()
	user.SetAge(25)
	CheckCount(t, 1)

	// 测试 SetBaseInfo 方法是否触发了 NotifyDirty
	ResetCount()
	baseInfo := dirty_out.NewBaseInfo()
	baseInfo.SetLv(10)
	baseInfo.SetExp(100)
	user.SetBaseInfo(baseInfo)
	CheckCount(t, 1)

	ResetCount()
	baseInfo.SetLv(1)
	CheckCount(t, 1)
}

func TestUserJsonMarshal(t *testing.T) {
	baseInfo := dirty_out.NewBaseInfo()
	baseInfo.SetLv(10)
	baseInfo.SetExp(100)
	b, err := json.Marshal(baseInfo)
	if err != nil {
		t.Error("error: ", err)
	}
	t.Log(string(b))

	if string(b) != `{"lv":10,"exp":100}` {
		t.Error("json marshal failed.")
	}
}

func TestUserJsonUnmarshal(t *testing.T) {
	baseInfo := dirty_out.NewBaseInfo()
	jsonStr := `{"lv":20,"exp":300}`
	err := json.Unmarshal([]byte(jsonStr), &baseInfo)
	if err != nil {
		t.Error("error:", err)
	}
	t.Logf("lv:%d, exp:%d\n", baseInfo.GetLv(), baseInfo.GetExp())

	if baseInfo.GetLv() != 20 || baseInfo.GetExp() != 300 {
		t.Error("json unmarshal failed.")
	}
}
