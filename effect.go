package nbattle

import "github.com/chompy/nbattle_go/internal/event"

type EffectDef struct {
	BaseObject
	Name string
	new  func() Effect
}

func (d *EffectDef) Serialize() []byte {
	return nil
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
