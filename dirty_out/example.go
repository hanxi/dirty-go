package dirty_out

/*
type Person struct {
	name    string
	age     int
	friends []*Person
	peoples map[string]*Person
}

type User struct {
	baseInfo *Person
	score    uint32
}
*/

// generate struct
type Person struct {
	Base
	name          string
	age           int
	_wrap_friends *WrapPersonFriends
	_wrap_peoples *WrapPersonPeoples
}

type WrapPersonFriends struct {
	Base
	friends []*Person
}

type WrapPersonPeoples struct {
	Base
	peoples map[string]*Person
}

type User struct {
	Base
	baseInfo *Person
	score    uint32
}

// generate func

func NewPerson() *Person {
	p := &Person{}
	p.self = p
	p.root = p
	return p
}

func (p *Person) SetName(value string) {
	if p == nil {
		return
	}
	p.name = value
	p.NotifyDirty()
}

func (p *Person) GetName() string {
	if p == nil {
		return ""
	}
	return p.name
}

func (p *Person) SetAge(value int) {
	if p == nil {
		return
	}
	p.age = value
	p.NotifyDirty()
}

func (p *Person) GetAge() int {
	if p == nil {
		return 0
	}
	return p.age
}

func (p *Person) SetFriends(value *WrapPersonFriends) {
	if p == nil {
		return
	}
	p._wrap_friends = value
	value.root = p.root
	for _, v := range value.friends {
		v.root = p.root
	}
	p.NotifyDirty()
}

func (p *Person) GetFriends() *WrapPersonFriends {
	if p == nil {
		return nil
	}
	return p._wrap_friends
}

func NewWrapPersonFriends() *WrapPersonFriends {
	p := &WrapPersonFriends{}
	p.friends = make([]*Person, 0)
	p.self = p
	p.root = p
	return p
}

func NewWrapPersonFriendsFromSlice(friends []*Person) *WrapPersonFriends {
	p := &WrapPersonFriends{}
	p.friends = make([]*Person, 0)
	p.friends = append(p.friends, friends...)
	p.self = p
	p.root = p
	return p
}

func (p *WrapPersonFriends) Append(value *Person) {
	if p == nil {
		return
	}
	p.friends = append(p.friends, value)
	value.root = p.root
	p.NotifyDirty()
}

func (p *WrapPersonFriends) Foreach(f func(*Person)) {
	if p == nil {
		return
	}
	for _, v := range p.friends {
		f(v)
	}
}

func (p *Person) SetPeoples(value *WrapPersonPeoples) {
	if p == nil {
		return
	}
	p._wrap_peoples = value
	value.root = p.root
	for _, v := range value.peoples {
		v.root = p.root
	}
	p.NotifyDirty()
}

func (p *Person) GetPeoples() *WrapPersonPeoples {
	if p == nil {
		return nil
	}
	return p._wrap_peoples
}

func NewWrapPersonPeoples() *WrapPersonPeoples {
	p := &WrapPersonPeoples{}
	p.peoples = make(map[string]*Person)
	p.self = p
	p.root = p
	return p
}

func NewWrapPersonPeoplesFromMap(peoples map[string]*Person) *WrapPersonPeoples {
	p := &WrapPersonPeoples{}
	p.peoples = make(map[string]*Person)
	for k, v := range peoples {
		p.peoples[k] = v
	}
	p.self = p
	p.root = p
	return p
}

func (p *WrapPersonPeoples) Get(key string) *Person {
	if p == nil {
		return nil
	}
	return p.peoples[key]
}

func (p *WrapPersonPeoples) Set(key string, value *Person) {
	if p == nil {
		return
	}
	p.peoples[key] = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *WrapPersonPeoples) Delete(key string) {
	if p == nil {
		return
	}
	delete(p.peoples, key)
	p.NotifyDirty()
}

func (p *WrapPersonPeoples) Foreach(f func(string, *Person)) {
	if p == nil {
		return
	}
	for k, v := range p.peoples {
		f(k, v)
	}
}

func NewUser() *User {
	p := &User{}
	p.self = p
	p.root = p
	return p
}

func (p *User) SetBaseInfo(value *Person) {
	if p == nil {
		return
	}
	p.baseInfo = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *User) GetBaseInfo() *Person {
	if p == nil {
		return nil
	}
	return p.baseInfo
}

func (p *User) SetScore(value uint32) {
	if p == nil {
		return
	}
	p.score = value
	p.NotifyDirty()
}

func (p *User) GetScore() uint32 {
	if p == nil {
		return 0
	}
	return p.score
}
