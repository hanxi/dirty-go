package dirty_out

type Person struct {
	Base
	name          string
	age           int
	_wrap_friends *ArrPersonFriends
	_wrap_peoples *MapPersonPeoples
}

type ArrPersonFriends struct {
	Base
	friends []*Person
}

func NewArrPersonFriends() *ArrPersonFriends {
	p := &ArrPersonFriends{}
	p.friends = make([]*Person, 0)
	p.self = p
	p.root = p
	return p
}

func NewArrPersonFriendsFromSlice(friends []*Person) *ArrPersonFriends {
	p := &ArrPersonFriends{}
	p.friends = make([]*Person, 0)
	p.friends = append(p.friends, friends...)
	p.self = p
	p.root = p
	return p
}

func (p *ArrPersonFriends) Append(value *Person) {
	if p == nil {
		return
	}
	p.friends = append(p.friends, value)
	value.root = p.root
	p.NotifyDirty()
}

func (p *ArrPersonFriends) Foreach(f func(*Person)) {
	if p == nil {
		return
	}
	for _, v := range p.friends {
		f(v)
	}
}

type MapPersonPeoples struct {
	Base
	peoples map[string]*Person
}

func NewMapPersonPeoples() *MapPersonPeoples {
	p := &MapPersonPeoples{}
	p.peoples = make(map[string]*Person, 0)
	p.self = p
	p.root = p
	return p
}

func NewMapPersonPeoplesFromMap(peoples map[string]*Person) *MapPersonPeoples {
	p := &MapPersonPeoples{}
	p.peoples = make(map[string]*Person)
	for k, v := range peoples {
		p.peoples[k] = v
	}
	p.self = p
	p.root = p
	return p
}

func (p *MapPersonPeoples) Get(key string) *Person {
	if p == nil {
		return nil
	}
	return p.peoples[key]
}

func (p *MapPersonPeoples) Set(key string, value *Person) {
	if p == nil {
		return
	}
	p.peoples[key] = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *MapPersonPeoples) Delete(key string) {
	if p == nil {
		return
	}
	delete(p.peoples, key)
	p.NotifyDirty()
}

func (p *MapPersonPeoples) Foreach(f func(string, *Person)) {
	if p == nil {
		return
	}
	for k, v := range p.peoples {
		f(k, v)
	}
}

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

func (p *Person) SetFriends(value *ArrPersonFriends) {
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

func (p *Person) GetFriends() *ArrPersonFriends {
	if p == nil {
		return nil
	}
	return p._wrap_friends
}

func (p *Person) SetPeoples(value *MapPersonPeoples) {
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

func (p *Person) GetPeoples() *MapPersonPeoples {
	if p == nil {
		return nil
	}
	return p._wrap_peoples
}

type Man struct {
	Base
	baseInfo *Person
	score    uint32
}

func NewMan() *Man {
	p := &Man{}
	p.self = p
	p.root = p
	return p
}

func (p *Man) SetBaseInfo(value *Person) {
	if p == nil {
		return
	}
	p.baseInfo = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *Man) GetBaseInfo() *Person {
	if p == nil {
		return nil
	}
	return p.baseInfo
}

func (p *Man) SetScore(value uint32) {
	if p == nil {
		return
	}
	p.score = value
	p.NotifyDirty()
}

func (p *Man) GetScore() uint32 {
	if p == nil {
		return 0
	}
	return p.score
}
