package nbattle

import (
	"slices"

	"github.com/chompy/nbattle_go/internal/event"
)

type CombatantEffect struct {
	Effect    Effect
	EffectCtx *EffectCtx
}

type Combatant struct {
	BaseObject
	active  bool
	stats   []*Stat
	effects []*CombatantEffect
}

func (c *Combatant) GetType() ObjectType {
	return ObjectTypeCombatant
}

func (c *Combatant) Reset() {
	for _, stat := range c.stats {
		stat.Reset()
	}
}

func (c *Combatant) IsActive() bool {
	return c.active
}

func (c *Combatant) SetActive(active bool) {
	if active != c.active {
		c.active = active
		c.ctx.EmitEvent(&event.CombatantUpdate{CombatantID: c.GetID(), Active: c.active})
	}
}

func (c *Combatant) GetStat(obj any) *Stat {
	def, ok := c.ctx.GetObject(obj).(*StatDef)
	if !ok {
		return nil
	}
	if c.stats == nil {
		c.stats = make([]*Stat, 0)
	}
	for _, stat := range c.stats {
		if stat.GetDef().GetID() == def.GetID() {
			return stat
		}
	}
	stat := &Stat{def, 0, nil}
	c.stats = append(c.stats, stat)
	return stat
}

func (c *Combatant) GetStats() []*Stat {
	return c.stats
}

func (c *Combatant) SetEffect(effectDefObj any, potency int, sourceObj any) error {
	effectDef, ok := c.ctx.GetObject(effectDefObj).(*EffectDef)
	if !ok {
		return ErrObjectNotFound
	}
	source := c.ctx.GetObject(sourceObj)

	if err := c.removeEffect(effectDef); err != nil {
		return err
	}
	if potency <= 0 {
		return nil
	}

	return c.addEffect(effectDef, potency, source)
}

func (c *Combatant) addEffect(effectDef *EffectDef, potency int, source Object) error {
	effect := effectDef.new()
	effectCtx := &EffectCtx{c.ctx, effectDef, potency, c, source}
	combatantEffect := &CombatantEffect{effect, effectCtx}
	c.effects = append(c.effects, combatantEffect)
	sourceID := 0
	if source != nil {
		sourceID = source.GetID()
	}
	c.ctx.addEffectToStack(effect)
	effect.OnAdd(effectCtx)
	c.ctx.removeEffectFromStack(effect)
	c.ctx.EmitEvent(&event.CombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID(), Potency: potency, SourceID: sourceID})
	return nil
}

func (c *Combatant) removeEffect(effectDef *EffectDef) error {
	c.effects = slices.DeleteFunc(c.effects, func(e *CombatantEffect) bool {
		if e.EffectCtx.Def.GetID() == effectDef.GetID() {
			c.ctx.addEffectToStack(e.Effect)
			e.Effect.OnRemove(e.EffectCtx)
			c.ctx.removeEffectFromStack(e.Effect)
			return true
		}
		return false
	})
	c.ctx.EmitEvent(&event.CombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID(), Potency: 0, SourceID: 0})
	return nil
}

func (c *Combatant) HasEffect(effectDefObj any) bool {
	effectDef, ok := c.ctx.GetObject(effectDefObj).(*EffectDef)
	if !ok {
		return false
	}
	return slices.ContainsFunc(c.effects, func(e *CombatantEffect) bool {
		return e.EffectCtx.Def.GetID() == effectDef.GetID()
	})
}

func (c *Combatant) HandleEffectEvent(event event.Event) error {
	for _, effect := range c.effects {
		if !c.ctx.isEffectInStack(effect.Effect) {
			c.ctx.addEffectToStack(effect.Effect)
			effect.Effect.OnEvent(effect.EffectCtx, event)
			c.ctx.removeEffectFromStack(effect.Effect)
		}
	}
	return nil
}
