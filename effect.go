package nbattle

type EffectDef struct {
	objectBase
	new func() Effect
}

func (d *EffectDef) Type() objectType {
	return objectTypeEffectDef
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
	OnEvent(ctx *EffectCtx, event *Event)
}
