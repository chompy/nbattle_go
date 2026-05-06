package nbattle

import (
	"slices"

	"github.com/chompy/nbattle_go/event"
)

type EffectDef struct {
	BaseObject
	name string
	new  func() Effect
}

func (d *EffectDef) GetType() ObjectType {
	return ObjectTypeEffectDef
}

func (d *EffectDef) GetName() string {
	return d.name
}

type EffectCtx struct {
	Ctx     *Context
	Def     *EffectDef
	Potency int
	Target  *Combatant
	Source  Object
}

func (e *EffectCtx) EmitTrigger(triggerDefObj any) error {
	triggerDef, err := e.Ctx.GetObject(triggerDefObj)
	if err != nil {
		return err
	}
	sourceID := 0
	if e.Source != nil {
		sourceID = e.Source.GetID()
	}
	e.Ctx.EmitEvent(&event.Trigger{
		TriggerDefID:   triggerDef.GetID(),
		EffectDefID:    e.Def.GetID(),
		EffectTargetID: e.Target.GetID(),
		EffectSourceID: sourceID,
		EffectPotency:  e.Potency,
	})
	return nil
}

type Effect interface {
	OnAdd(ctx *EffectCtx)
	OnRemove(ctx *EffectCtx)
	OnEvent(ctx *EffectCtx, event event.Event)
}

func (c *Context) addEffectToStack(effect Effect) {
	c.effectStack = append(c.effectStack, effect)
}

func (c *Context) removeEffectFromStack(effect Effect) {
	c.effectStack = slices.DeleteFunc(c.effectStack, func(e Effect) bool { return e == effect })
}

func (c *Context) isEffectInStack(effect Effect) bool {
	return slices.Contains(c.effectStack, effect)
}
