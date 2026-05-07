package nbattle

import (
	"slices"

	"github.com/chompy/nbattle_go/event"
)

// CombatantEffect pairs an effect instance to an effect context.
type CombatantEffect struct {
	Effect    Effect
	EffectCtx *EffectCtx
}

// Combatant is an entity that can interact with other combatants.
type Combatant struct {
	BaseObject
	active  bool
	stats   []*Stat
	effects []*CombatantEffect
	flags   uint64
}

func (c *Combatant) GetType() ObjectType {
	return ObjectTypeCombatant
}

// Reset resets the combatant's stats to their base values.
func (c *Combatant) Reset() {
	for _, stat := range c.stats {
		stat.Reset()
	}
}

// IsActive returns true if the combatant is active.
func (c *Combatant) IsActive() bool {
	return c.active
}

// SetActive sets the combatant's active state.
func (c *Combatant) SetActive(active bool) {
	if active != c.active {
		c.active = active
		c.emitUpdateEvent()
	}
}

// GetStat retrieves the combatant's stat from a stat definition object.
func (c *Combatant) GetStat(statDefObj any) (*Stat, error) {
	statDefObj, err := c.ctx.GetObject(statDefObj)
	if err != nil {
		return nil, err
	}
	statDef, ok := statDefObj.(*StatDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	if c.stats == nil {
		c.stats = make([]*Stat, 0)
	}
	for _, stat := range c.stats {
		if stat.GetDef().GetID() == statDef.GetID() {
			return stat, nil
		}
	}
	stat := &Stat{statDef, 0, nil}
	c.stats = append(c.stats, stat)
	return stat, nil
}

// GetStats retrieves all the combatant's stats.
func (c *Combatant) GetStats() []*Stat {
	return c.stats
}

// SetEffect adds an effect to the combatant. Potency is the power of the effect.
// A potency of zero removes the effect.
func (c *Combatant) SetEffect(effectDefObj any, potency int, sourceObj any) error {
	effectDefObj, err := c.ctx.GetObject(effectDefObj)
	if err != nil {
		return err
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return ErrUnexpectedObjectType
	}
	source, _ := c.ctx.GetObject(sourceObj)
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

// HasEffect returns true if the combatant has the given effect.
func (c *Combatant) HasEffect(effectDefObj any) bool {
	effectDefObj, err := c.ctx.GetObject(effectDefObj)
	if err != nil {
		return false
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return false
	}
	return slices.ContainsFunc(c.effects, func(e *CombatantEffect) bool {
		return e.EffectCtx.Def.GetID() == effectDef.GetID()
	})
}

// SetFlag sets the given flag to true or false.
func (c *Combatant) SetFlag(flag any, on bool) {
	var flagValue uint64
	switch f := flag.(type) {
	case string:
		flagValue = c.ctx.GetFlagByName(f)
	case uint64:
		flagValue = f
	case int64:
		flagValue = uint64(f)
	case int:
		flagValue = uint64(f)
	case float64:
		flagValue = uint64(f)
	case float32:
		flagValue = uint64(f)
	default:
		return
	}
	if on {
		c.flags |= flagValue
	} else {
		c.flags &^= flagValue
	}
	c.emitUpdateEvent()
}

// HasFlag returns true if the given flag is set.
func (c *Combatant) HasFlag(flag any) bool {
	var flagValue uint64
	switch f := flag.(type) {
	case string:
		flagValue = c.ctx.GetFlagByName(f)
	case uint64:
		flagValue = f
	case int64:
		flagValue = uint64(f)
	case int:
		flagValue = uint64(f)
	case float64:
		flagValue = uint64(f)
	case float32:
		flagValue = uint64(f)
	default:
		return false
	}
	return (c.flags & flagValue) != 0
}

// GetFlags returns the current flags value.
func (c *Combatant) GetFlags() uint64 {
	return c.flags
}

// HandleEffectEvent sends the given event to all the effects of this combatant.
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

func (c *Combatant) emitUpdateEvent() {
	c.ctx.EmitEvent(&event.CombatantUpdate{CombatantID: c.GetID(), Active: c.active, Flags: c.flags})
}
