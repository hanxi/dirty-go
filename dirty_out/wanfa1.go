package dirty_out

type Wanfa1 struct {
	Base
	name string
}

func NewWanfa1() *Wanfa1 {
	p := &Wanfa1{}
	p.self = p
	p.root = p
	return p
}

func (p *Wanfa1) SetName(value string) {
	if p == nil {
		return
	}
	p.name = value
	p.NotifyDirty()
}

func (p *Wanfa1) GetName() string {
	if p == nil {
		return ""
	}
	return p.name
}
