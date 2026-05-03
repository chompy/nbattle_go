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
	Team    int
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
	stat := &Stat{c.ctx.newObject(), def, 0, nil}
	c.ctx.EmitEvent(&event.AddCombatantStat{CombatantID: c.GetID(), StatID: stat.GetID(), StatDefID: def.GetID()})
	c.ctx.objects = append(c.ctx.objects, stat)
	c.stats = append(c.stats, stat)
	return stat
}

func (c *Combatant) GetStats() []*Stat {
	return c.stats
}

func (c *Combatant) AddEffect(effectDefObj any, sourceObj any) error {
	effectDef, ok := c.ctx.GetObject(effectDefObj).(*EffectDef)
	if !ok {
		return ErrObjectNotFound
	}
	effect := effectDef.new()
	effectCtx := &EffectCtx{c.ctx, effectDef, nil, c}
	combatantEffect := &CombatantEffect{effect, effectCtx}
	c.effects = append(c.effects, combatantEffect)
	effect.OnAdd(effectCtx)
	sourceID := 0
	if sourceObj != nil {
		if src, ok2 := c.ctx.GetObject(sourceObj).(*Combatant); ok2 {
			sourceID = src.GetID()
			effectCtx.Source = src
		}
	}
	c.ctx.EmitEvent(&event.AddCombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID(), SourceID: sourceID})
	return nil
}

func (c *Combatant) RemoveEffect(effectDefObj any) error {
	effectDef, ok := c.ctx.GetObject(effectDefObj).(*EffectDef)
	if !ok {
		return ErrObjectNotFound
	}
	c.effects = slices.DeleteFunc(c.effects, func(e *CombatantEffect) bool {
		if e.EffectCtx.Def.GetID() == effectDef.GetID() {
			e.Effect.OnRemove(e.EffectCtx)
			return true
		}
		return false
	})
	c.ctx.EmitEvent(&event.RemoveCombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID()})
	return nil
}

func (c *Combatant) HandleEffectEvent(event event.Event) error {
	for _, effect := range c.effects {
		effect.Effect.OnEvent(effect.EffectCtx, event)
	}
	return nil
}
