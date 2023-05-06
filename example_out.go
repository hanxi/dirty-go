package main

/*
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

func (p *Person) SetFriends(value []*Person) {
	if p == nil {
		return
	}
	p.friends = value
	for _, v := range value {
		v.root = p.root
	}
	p.NotifyDirty()
}

func (p *Person) GetFriends() []*Person {
	if p == nil {
		return nil
	}
	return p.friends
}

func (p *Person) AppendFriends(value *Person) {
	if p == nil {
		return
	}
	p.friends = append(p.friends, value)
	value.root = p.root
	p.NotifyDirty()
}

func (p *Person) SetPeoples(value map[string]*Person) {
	if p == nil {
		return
	}
	p.peoples = value
	for _, v := range value {
		v.root = p.root
	}
	p.NotifyDirty()
}

func (p *Person) GetPeoples() map[string]*Person {
	if p == nil {
		return nil
	}
	return p.peoples
}

func (p *Person) PutPeoples(key string, value *Person) {
	if p == nil {
		return
	}
	p.peoples[key] = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *Person) LookupPeoples(key string) *Person {
	if p == nil {
		return nil
	}
	return p.peoples[key]
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
*/
