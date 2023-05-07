// Code generated by dirty-go; DO NOT EDIT.
package dirty_out

type BaseInfo struct {
	Base
	lv uint32
	exp uint32
}

func NewBaseInfo() *BaseInfo {
	p := &BaseInfo{}
	p.self = p
	p.root = p
	return p
}

type Resource struct {
	Base
	id uint32
	value uint32
	size uint32
}

func NewResource() *Resource {
	p := &Resource{}
	p.self = p
	p.root = p
	return p
}

type Friend struct {
	Base
	uid uint32
	name string
}

func NewFriend() *Friend {
	p := &Friend{}
	p.self = p
	p.root = p
	return p
}

type User struct {
	Base
	name string
	age int
	baseInfo *BaseInfo
	_map_resources *MapUserResources
	_arr_friends *ArrUserFriends
}

func NewUser() *User {
	p := &User{}
	p.self = p
	p.root = p
	return p
}
