package combat

import (
	"slices"

	"github.com/chompy/nbattle_go/internal/event"
)

// EffectDef is a definition of an effect.
// An effect can be applied to combatants and used to dynamically alter the combatant's stats.
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

// EffectCtx is the context for an effect that has been applied to a specific combatant.
type EffectCtx struct {
	Ctx     *Context
	Def     *EffectDef
	Potency int
	Target  *Combatant
	Source  Object
}

// EmitTrigger emits a trigger event for the given trigger definition.
func (e *EffectCtx) EmitTrigger(triggerDefObj any) error {
	triggerDefObj, err := e.Ctx.GetObject(triggerDefObj)
	if err != nil {
		return err
	}
	triggerDef, ok := triggerDefObj.(*TriggerDef)
	if !ok {
		return ErrUnexpectedObjectType
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

// Effect is an interface for effects.
type Effect interface {
	OnAdd(ctx *EffectCtx)                      // OnAdd is called when the effect is first applied to a combatant.
	OnRemove(ctx *EffectCtx)                   // OnRemove is called when the effect is removed from a combatant.
	OnEvent(ctx *EffectCtx, event event.Event) // OnEvent is called whenever an outside event is emitted.
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
