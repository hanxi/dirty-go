package main

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
	parent   DataObject
	root     DataObject
	self     DataObject
}

func (x *Base) NotifyDirty() {
	if x.observer != nil {
		x.observer.OnDirty(x)
	}
	if x.root != nil && x.root != x.self {
		// 非根节点往上传递消息
		x.root.NotifyDirty()
	}
}

func (x *Base) Attach(o Observer) {
	x.observer = o
}
