package nbattle

type ObjectType int

const (
	ObjectTypeStatDef = iota
	ObjectTypeEffectDef
	ObjectTypeStat
	ObjectTypeCombatant
	ObjectTypeUnknown
	ObjectTypeError
)

type Object interface {
	GetID() int
	GetType() ObjectType
}

type BaseObject struct {
	id  int
	ctx *Context
}

func (o *BaseObject) GetID() int {
	return o.id
}
