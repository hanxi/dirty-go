package dirty_out

type BaseInfo struct {
	Base
	lv  uint32
	exp uint32
}

func NewBaseInfo() *BaseInfo {
	p := &BaseInfo{}
	p.self = p
	p.root = p
	return p
}

func (p *BaseInfo) SetLv(value uint32) {
	if p == nil {
		return
	}
	p.lv = value
	p.NotifyDirty()
}

func (p *BaseInfo) GetLv() uint32 {
	if p == nil {
		return 0
	}
	return p.lv
}

func (p *BaseInfo) SetExp(value uint32) {
	if p == nil {
		return
	}
	p.exp = value
	p.NotifyDirty()
}

func (p *BaseInfo) GetExp() uint32 {
	if p == nil {
		return 0
	}
	return p.exp
}

type Resource struct {
	Base
	id    uint32
	value uint32
	size  uint32
}

func NewResource() *Resource {
	p := &Resource{}
	p.self = p
	p.root = p
	return p
}

func (p *Resource) SetId(value uint32) {
	if p == nil {
		return
	}
	p.id = value
	p.NotifyDirty()
}

func (p *Resource) GetId() uint32 {
	if p == nil {
		return 0
	}
	return p.id
}

func (p *Resource) SetValue(value uint32) {
	if p == nil {
		return
	}
	p.value = value
	p.NotifyDirty()
}

func (p *Resource) GetValue() uint32 {
	if p == nil {
		return 0
	}
	return p.value
}

func (p *Resource) SetSize(value uint32) {
	if p == nil {
		return
	}
	p.size = value
	p.NotifyDirty()
}

func (p *Resource) GetSize() uint32 {
	if p == nil {
		return 0
	}
	return p.size
}

type Friend struct {
	Base
	uid  uint32
	name string
}

func NewFriend() *Friend {
	p := &Friend{}
	p.self = p
	p.root = p
	return p
}

func (p *Friend) SetUid(value uint32) {
	if p == nil {
		return
	}
	p.uid = value
	p.NotifyDirty()
}

func (p *Friend) GetUid() uint32 {
	if p == nil {
		return 0
	}
	return p.uid
}

func (p *Friend) SetName(value string) {
	if p == nil {
		return
	}
	p.name = value
	p.NotifyDirty()
}

func (p *Friend) GetName() string {
	if p == nil {
		return ""
	}
	return p.name
}

type User struct {
	Base
	name            string
	age             int
	baseInfo        *BaseInfo
	_wrap_resources *MapUserResources
	_wrap_friends   *ArrUserFriends
}

type MapUserResources struct {
	Base
	resources map[uint32]*Resource
}

func NewMapUserResources() *MapUserResources {
	p := &MapUserResources{}
	p.resources = make(map[uint32]*Resource, 0)
	p.self = p
	p.root = p
	return p
}

func NewMapUserResourcesFromMap(resources map[uint32]*Resource) *MapUserResources {
	p := &MapUserResources{}
	p.resources = make(map[uint32]*Resource)
	for k, v := range resources {
		p.resources[k] = v
	}
	p.self = p
	p.root = p
	return p
}

func (p *MapUserResources) Get(key uint32) *Resource {
	if p == nil {
		return nil
	}
	return p.resources[key]
}

func (p *MapUserResources) Set(key uint32, value *Resource) {
	if p == nil {
		return
	}
	p.resources[key] = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *MapUserResources) Delete(key uint32) {
	if p == nil {
		return
	}
	delete(p.resources, key)
	p.NotifyDirty()
}

func (p *MapUserResources) Foreach(f func(uint32, *Resource)) {
	if p == nil {
		return
	}
	for k, v := range p.resources {
		f(k, v)
	}
}

type ArrUserFriends struct {
	Base
	friends []*Friend
}

func NewArrUserFriends() *ArrUserFriends {
	p := &ArrUserFriends{}
	p.friends = make([]*Friend, 0)
	p.self = p
	p.root = p
	return p
}

func NewArrUserFriendsFromSlice(friends []*Friend) *ArrUserFriends {
	p := &ArrUserFriends{}
	p.friends = make([]*Friend, 0)
	p.friends = append(p.friends, friends...)
	p.self = p
	p.root = p
	return p
}

func (p *ArrUserFriends) Append(value *Friend) {
	if p == nil {
		return
	}
	p.friends = append(p.friends, value)
	value.root = p.root
	p.NotifyDirty()
}

func (p *ArrUserFriends) Foreach(f func(*Friend)) {
	if p == nil {
		return
	}
	for _, v := range p.friends {
		f(v)
	}
}

func NewUser() *User {
	p := &User{}
	p.self = p
	p.root = p
	return p
}

func (p *User) SetName(value string) {
	if p == nil {
		return
	}
	p.name = value
	p.NotifyDirty()
}

func (p *User) GetName() string {
	if p == nil {
		return ""
	}
	return p.name
}

func (p *User) SetAge(value int) {
	if p == nil {
		return
	}
	p.age = value
	p.NotifyDirty()
}

func (p *User) GetAge() int {
	if p == nil {
		return 0
	}
	return p.age
}

func (p *User) SetBaseInfo(value *BaseInfo) {
	if p == nil {
		return
	}
	p.baseInfo = value
	value.root = p.root
	p.NotifyDirty()
}

func (p *User) GetBaseInfo() *BaseInfo {
	if p == nil {
		return nil
	}
	return p.baseInfo
}

func (p *User) SetResources(value *MapUserResources) {
	if p == nil {
		return
	}
	p._wrap_resources = value
	value.root = p.root
	for _, v := range value.resources {
		v.root = p.root
	}
	p.NotifyDirty()
}

func (p *User) GetResources() *MapUserResources {
	if p == nil {
		return nil
	}
	return p._wrap_resources
}

func (p *User) SetFriends(value *ArrUserFriends) {
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

func (p *User) GetFriends() *ArrUserFriends {
	if p == nil {
		return nil
	}
	return p._wrap_friends
}
