package dirty_out

type Observer interface {
	OnDirty(interface{})
}

type DataObject interface {
	NotifyDirty()
	Attach(o Observer)
}

type Base struct {
	DataObject
	observer Observer
	root     DataObject
	self     DataObject
}

func (x *Base) NotifyDirty() {
	if x.observer != nil {
		x.observer.OnDirty(x)
	}
	if x.root != nil && x.root != x.self {
		x.root.NotifyDirty()
	}
}

func (x *Base) Attach(o Observer) {
	x.observer = o
}
