package nbattle

type objectType int

const (
	objectTypeStatDef objectType = iota
	objectTypeStat
	objectTypeCombatant
	objectTypeEffectDef
	objectTypeEffect
)

type object interface {
	ID() int
	Type() objectType
	Serialize() []byte
}

type objectBase struct {
	id  int
	ctx *Context
}

func (o *objectBase) ID() int {
	return o.id
}
