package nbattle

import "slices"

type combatantEffect struct {
	Effect    Effect
	EffectCtx *EffectCtx
}

type Combatant struct {
	objectBase
	stats   []*Stat
	effects []*combatantEffect
}

func (d *Combatant) Type() objectType {
	return objectTypeCombatant
}

func (d *Combatant) Serialize() []byte {
	return nil
}

func (c *Combatant) Reset() {
	for _, stat := range c.stats {
		stat.Reset()
	}
}

func (c *Combatant) Stat(def *StatDef) *Stat {
	for _, stat := range c.stats {
		if stat.Def().ID() == def.ID() {
			return stat
		}
	}
	stat := &Stat{c.ctx.newObject(), def, 0, nil}
	c.ctx.EmitEvent(EventTypeCombatantStatAdd, c.ID(), stat.ID(), def.ID())
	c.ctx.objects = append(c.ctx.objects, stat)
	c.stats = append(c.stats, stat)
	return stat
}

func (c *Combatant) AddEffect(effectDef *EffectDef, source *Combatant) {
	effect := effectDef.new()
	effectCtx := &EffectCtx{c.ctx, effectDef, source, c}
	combatantEffect := &combatantEffect{effect, effectCtx}
	c.effects = append(c.effects, combatantEffect)
	effect.OnAdd(effectCtx)
	c.ctx.EmitEvent(EventTypeCombatantEffectAdd, c.ID(), effectDef.ID(), source.ID())
}

func (c *Combatant) RemoveEffect(effectDef *EffectDef) {
	c.effects = slices.DeleteFunc(c.effects, func(e *combatantEffect) bool {
		if e.EffectCtx.Def.ID() == effectDef.ID() {
			e.Effect.OnRemove(e.EffectCtx)
			return true
		}
		return false
	})
}

func (c *Combatant) HandleEffectEvent(event *Event) error {
	for _, effect := range c.effects {
		effect.Effect.OnEvent(effect.EffectCtx, event)
	}
	return nil
}
