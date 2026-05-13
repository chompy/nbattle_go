package nbattle

type ObjectType int

const (
	ObjectTypeStatDef = iota
	ObjectTypeEffectDef
	ObjectTypeTriggerDef
	ObjectTypeCombatant
	ObjectTypeUnknown
	ObjectTypeError
)

// Object is an interface for all objects managed by NBattle.
type Object interface {
	GetID() int          // Get the object's unique ID.
	GetType() ObjectType // Get the object's type.
	String() string      // Get the string representation of the object.
}

// BaseObject is a base struct for all objects.
type BaseObject struct {
	id  int
	ctx *Context
}

// Retrieve the object ID.
func (o *BaseObject) GetID() int {
	return o.id
}
