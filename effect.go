package nbattle

import (
	"slices"

	"github.com/chompy/nbattle_go/internal/event"
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
	Ctx    *Context
	Def    *EffectDef
	Source *Combatant
	Target *Combatant
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
