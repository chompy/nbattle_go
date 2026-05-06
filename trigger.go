package nbattle

import "github.com/chompy/nbattle_go/event"

type TriggerDef struct {
	BaseObject
	name string
}

func (d *TriggerDef) GetType() ObjectType {
	return ObjectTypeEffectDef
}

func (d *TriggerDef) GetName() string {
	return d.name
}

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
