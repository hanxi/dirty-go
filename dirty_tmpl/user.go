package dirty_tmpl

type BaseInfo struct {
	lv  uint32
	exp uint32
}

type Resource struct {
	id    uint32
	value uint32
	size  uint32
}

type Friend struct {
	uid  uint32
	name string
}

type User struct {
	name      string
	age       int
	baseInfo  *BaseInfo
	resources map[uint32]*Resource
	friends   []*Friend
}
