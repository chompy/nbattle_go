package nbattle

import "fmt"

type ObjectType int

const (
	ObjectTypeStatDef = iota
	ObjectTypeEffectDef
	ObjectTypeTriggerDef
	ObjectTypeCombatant
	ObjectTypeUnknown
	ObjectTypeError
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypeStatDef:
		return "stat_def"
	case ObjectTypeEffectDef:
		return "effect_def"
	case ObjectTypeTriggerDef:
		return "trigger_def"
	case ObjectTypeCombatant:
		return "combatant"
	case ObjectTypeUnknown:
		return "unknown"
	case ObjectTypeError:
		return "error"
	default:
		return fmt.Sprintf("unknown %d", o)
	}
}

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
