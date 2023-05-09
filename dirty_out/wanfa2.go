package dirty_out

type Wanfa2 struct {
	Base
	name string
}

func NewWanfa2() *Wanfa2 {
	p := &Wanfa2{}
	p.self = p
	p.root = p
	return p
}

func (p *Wanfa2) SetName(value string) {
	if p == nil {
		return
	}
	p.name = value
	p.NotifyDirty()
}

func (p *Wanfa2) GetName() string {
	if p == nil {
		return ""
	}
	return p.name
}
