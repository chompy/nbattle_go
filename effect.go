package nbattle

import (
	"fmt"
	"slices"

	"github.com/chompy/nbattle_go/event"
)

// EffectDef is a definition of an effect.
// An effect can be applied to combatants and used to dynamically alter the combatant's stats.
type EffectDef struct {
	BaseObject
	name string
	new  func() Effect
}

func (e *EffectDef) String() string {
	return fmt.Sprintf("<EffectDef name=%s id=%d>", e.name, e.GetID())
}

func (d *EffectDef) GetType() ObjectType {
	return ObjectTypeEffectDef
}

func (d *EffectDef) GetName() string {
	return d.name
}

// EffectCtx is the context for an effect that has been applied to a specific combatant.
type EffectContext struct {
	Def     *EffectDef
	Target  *Combatant
	Source  Object
	Potency int
}

func (e *EffectContext) String() string {
	return fmt.Sprintf("<EffectCtx def=%s target=%s source=%s potency=%d>", e.Def.String(), e.Target.String(), e.Source.String(), e.Potency)
}

// Effect is an interface for effects.
type Effect interface {
	OnAdd(ctx *Context, effect *EffectContext)                      // OnAdd is called when the effect is first applied to a combatant.
	OnRemove(ctx *Context, effect *EffectContext)                   // OnRemove is called when the effect is removed from a combatant.
	OnEvent(ctx *Context, effect *EffectContext, event event.Event) // OnEvent is called whenever an outside event is emitted.
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

func getCombatantEffectContext(target *Combatant, source Object, combatantEffect *combatantEffect) *EffectContext {
	potency := combatantEffect.sources[source.GetID()]
	return &EffectContext{
		Def:     combatantEffect.def,
		Potency: potency,
		Target:  target,
		Source:  source,
	}
}
