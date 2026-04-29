package nbattle

type Combatant struct {
	objectBase
	stats   []*Stat
	effects map[*EffectDef]Effect
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

func (c *Combatant) ApplyEffect(effectDef *EffectDef, source *Combatant) Effect {
	effect := effectDef.create(c, source)
	c.effects[effectDef] = effect
	effect.OnApply()
	c.ctx.EmitEvent(EventTypeCombatantEffectApply, c.ID(), effectDef.ID(), source.ID())
	return effect
}

func (c *Combatant) RemoveEffect(effectDef *EffectDef) {
	if c.effects[effectDef] != nil {
		c.effects[effectDef].OnRemove()
		delete(c.effects, effectDef)
		c.ctx.EmitEvent(EventTypeCombatantEffectRemove, c.ID(), effectDef.ID())
	}
}
