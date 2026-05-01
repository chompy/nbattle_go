package nbattle

type Object interface {
	GetID() int
}

type BaseObject struct {
	id  int
	ctx *Context
}

func (o *BaseObject) GetID() int {
	return o.id
}
