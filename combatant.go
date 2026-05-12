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
	// find the effect definition
	effectDefObj, err := c.ctx.GetObject(effectDefObj)
	if err != nil {
		return err
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return ErrUnexpectedObjectType
	}

	// find source object id if exists
	source, _ := c.ctx.GetObject(sourceObj)
	sourceID := 0
	if source != nil {
		sourceID = source.GetID()
	}

	// check if effect is already applied to the combatant
	var combatantEffect *CombatantEffect
	for _, ce := range c.effects {
		if ce.EffectCtx.Def.GetID() == effectDef.GetID() {
			combatantEffect = ce
			break
		}
	}

	// zero potency means effect should be removed
	if potency <= 0 && combatantEffect != nil {
		combatantEffect.EffectCtx.Potency = 0
		return nil
	}

	// no change, same potency as existing
	if combatantEffect != nil && potency == combatantEffect.EffectCtx.Potency {
		return nil
	}

	// create new instance of effect if not already applied
	if combatantEffect == nil {
		effect := effectDef.new()
		effectCtx := &EffectCtx{c.ctx, effectDef, 0, c, source}
		combatantEffect = &CombatantEffect{effect, effectCtx}
		c.effects = append(c.effects, combatantEffect)
	}
	combatantEffect.EffectCtx.Potency = potency

	// call effect's OnAdd function and emit event
	c.ctx.addEffectToStack(combatantEffect.Effect)
	combatantEffect.Effect.OnAdd(combatantEffect.EffectCtx)
	c.ctx.removeEffectFromStack(combatantEffect.Effect)
	return c.ctx.EmitEvent(&event.CombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID(), Potency: potency, SourceID: sourceID})
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

// HandleEffectEvent sends the given event to all the active effects of this combatant.
func (c *Combatant) HandleEffectEvent(event event.Event) error {
	for _, effect := range c.effects {
		if !c.ctx.isEffectInStack(effect.Effect) && effect.EffectCtx.Potency > 0 {
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

func (c *Combatant) processEffectRemovals() {
	c.effects = slices.DeleteFunc(c.effects, func(e *CombatantEffect) bool {
		if e.EffectCtx.Potency <= 0 {
			c.ctx.addEffectToStack(e.Effect)
			e.Effect.OnRemove(e.EffectCtx)
			c.ctx.removeEffectFromStack(e.Effect)
			c.ctx.EmitEvent(&event.CombatantEffect{TargetID: c.GetID(), EffectDefID: e.EffectCtx.Def.GetID(), Potency: 0, SourceID: 0})
			return true
		}
		return false
	})
}
