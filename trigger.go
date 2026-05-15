package nbattle

import (
	"fmt"
)

// TriggerDef is a definition of a custom event that can be triggered by an effect.
type TriggerDef struct {
	BaseObject
	name string
}

func (s *TriggerDef) String() string {
	return fmt.Sprintf("<TriggerDef name=%s id=%d>", s.name, s.GetID())
}

func (d *TriggerDef) GetType() ObjectType {
	return ObjectTypeTriggerDef
}

func (d *TriggerDef) GetName() string {
	return d.name
}
