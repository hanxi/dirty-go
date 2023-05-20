package dirty_tmpl

type BaseInfo struct {
	Lv  uint32 `json:"lv"`
	Exp uint32 `json:"exp"`
}

type Resource struct {
	Id    uint32
	Value uint32
	Size  uint32
}

type Friend struct {
	Uid  uint32
	Name string
}

type User struct {
	Name      string
	Age       int
	BaseInfo  *BaseInfo
	Resources map[uint32]*Resource
	Friends   []*Friend
}
