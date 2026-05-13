package nbattle

import (
	"fmt"

	"github.com/chompy/nbattle_go/event"
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

// EmitEvent emits a trigger event.
func (d *TriggerDef) EmitEvent(effectCtx *EffectCtx) {
	sourceID := 0
	if effectCtx.Source != nil {
		sourceID = effectCtx.Source.GetID()
	}
	d.ctx.EmitEvent(&event.Trigger{
		TriggerDefID:   d.GetID(),
		EffectDefID:    effectCtx.Def.GetID(),
		EffectTargetID: effectCtx.Target.GetID(),
		EffectSourceID: sourceID,
		EffectPotency:  effectCtx.Potency,
	})
}
