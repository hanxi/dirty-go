package dirty_tmpl

type BaseInfo struct {
	Lv  uint32 `json:"lv"`
	Exp uint32 `json:"exp"`
}

type Resource struct {
	Id    uint32 `json:"id"`
	Value uint32 `json:"value"`
	Size  uint32 `json:"size"`
}

type Friend struct {
	Uid  uint32 `json:"uid"`
	Name string `json:"name"`
}

type User struct {
	Name      string               `json:"name"`
	Age       int                  `json:"age"`
	Info      *BaseInfo            `json:"baseinfo"`
	Resources map[uint32]*Resource `json:resources`
	Friends   []*Friend            `json:friends`
}
