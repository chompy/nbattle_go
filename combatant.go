package nbattle

import (
	"fmt"
	"slices"

	"github.com/chompy/nbattle_go/event"
)

type combatantEffect struct {
	def      *EffectDef
	instance Effect
	sources  map[int]int
}

// Combatant is an entity that can interact with other combatants.
type Combatant struct {
	BaseObject
	active  bool
	stats   []*Stat
	effects []*combatantEffect
	flags   uint64
}

func (c *Combatant) String() string {
	return fmt.Sprintf("<Combatant id=%d>", c.GetID())
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
	statDef, err := c.ctx.GetStatDef(statDefObj)
	if err != nil {
		return nil, err
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

// SetEffect adds an effect to the combatant. The same effect can be applied
// multiple times from different sources. Potency is the power of the effect.
// A potency of zero removes the effect.
func (c *Combatant) SetEffect(effectDefObj any, sourceObj any, potency int) error {
	// find the effect definition
	effectDef, err := c.ctx.GetEffectDef(effectDefObj)
	if err != nil {
		return err
	}

	// find the effect source
	source, err := c.ctx.GetObject(sourceObj)
	if err != nil {
		return err
	}

	c.ctx.log.Debug("Set effect to combatant.", "effectDef", effectDef, "target", c, "source", source, "potency", potency)

	// check if effect is already applied to the combatant
	var combatantEft *combatantEffect
	for _, ce := range c.effects {
		if ce.def.GetID() == effectDef.GetID() {
			combatantEft = ce
			break
		}
	}

	// zero potency means effect should be removed
	if potency <= 0 {
		if combatantEft != nil && combatantEft.sources[source.GetID()] > 0 {
			c.ctx.addEffectToStack(combatantEft.instance)
			effectCtx := getCombatantEffectContext(c, source, combatantEft)
			combatantEft.instance.OnRemove(c.ctx, effectCtx)
			c.ctx.removeEffectFromStack(combatantEft.instance)

			combatantEft.sources[source.GetID()] = 0
		}
		return nil
	}

	// no change, same potency as existing
	if combatantEft != nil && potency == combatantEft.sources[source.GetID()] {
		return nil
	}

	// create new instance of effect if not already applied
	if combatantEft == nil {
		effect := effectDef.new()
		combatantEft = &combatantEffect{effectDef, effect, make(map[int]int)}
		c.effects = append(c.effects, combatantEft)
	}

	combatantEft.sources[source.GetID()] = potency

	// call effect's OnAdd function and emit event
	c.ctx.addEffectToStack(combatantEft.instance)
	effectCtx := getCombatantEffectContext(c, source, combatantEft)
	combatantEft.instance.OnAdd(c.ctx, effectCtx)
	c.ctx.removeEffectFromStack(combatantEft.instance)

	return c.ctx.EmitEvent(&event.CombatantEffect{TargetID: c.GetID(), EffectDefID: effectDef.GetID(), Potency: potency, SourceID: source.GetID()})
}

// HasEffect returns true if the combatant has the given effect active.
func (c *Combatant) HasEffect(effectDefObj any) bool {
	effectDef, err := c.ctx.GetEffectDef(effectDefObj)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(c.effects, func(e *combatantEffect) bool {
		if e.def.GetID() == effectDef.GetID() {
			for _, potency := range e.sources {
				if potency > 0 {
					return true
				}
			}
		}
		return false
	})
}

func (c *Combatant) emitUpdateEvent() {
	c.ctx.EmitEvent(&event.CombatantUpdate{CombatantID: c.GetID(), Active: c.active, Flags: c.flags})
}

func (c *Combatant) processEvent(event event.Event) error {
	for _, effect := range c.effects {
		if !c.ctx.isEffectInStack(effect.instance) {
			c.ctx.addEffectToStack(effect.instance)
			for sourceID, potency := range effect.sources {
				if potency > 0 {
					source, err := c.ctx.GetObject(sourceID)
					if err != nil {
						return err
					}
					effectCtx := getCombatantEffectContext(c, source, effect)
					effect.instance.OnEvent(c.ctx, effectCtx, event)
				}
			}
			c.ctx.removeEffectFromStack(effect.instance)
		}
	}
	return nil
}
