package nbattle

import "github.com/chompy/nbattle_go/internal/event"

// StatDef is a definition of a stat that can be applied to a combatant.
type StatDef struct {
	BaseObject
	name string
	min  int
	max  int
}

func (s *StatDef) GetType() ObjectType {
	return ObjectTypeStatDef
}

func (d *StatDef) GetName() string {
	return d.name
}

// GetMin is the minimum value of the stat.
func (d *StatDef) GetMin() int {
	return d.min
}

// GetMax is the maximum value of the stat.
func (d *StatDef) GetMax() int {
	return d.max
}

// Stat represents a single integer value that can be applied to a combatant.
// It represents the combatant's combat ability in one specific way (such as strength or intelegence).
// A stat can have mods applied to it that are temporary alterations applied by effects.
// Mods are a mapping of the effect definition ID that applied it and the ammount to offset the base stat value by.
type Stat struct {
	def  *StatDef
	base int
	mods map[int]int
}

// GetDef retrieves the stat's stat definition.
func (s *Stat) GetDef() *StatDef {
	return s.def
}

// GetBase retrieves the base value of the stat without any mods applied.
func (s *Stat) GetBase() int {
	return s.base
}

// SetBase sets the base value of the stat.
func (s *Stat) SetBase(value int) {
	statDef := s.GetDef()
	combatant, err := s.def.ctx.GetCombatantWithStat(s)
	if err == nil && combatant != nil {
		// fire combatant stat event, allow hooks to modify the stat value change
		evt := &event.CombatantStatBase{CombatantID: combatant.GetID(), StatDefID: statDef.GetID(), Value: value}
		statDef.ctx.EmitEvent(evt)
		value = evt.Value
	}
	s.base = min(statDef.GetMax(), max(statDef.GetMin(), value))
}

// AddBase adds the given amount to the base stat.
func (s *Stat) AddBase(value int) {
	s.SetBase(s.base + value)
}

// SubBase subtracts the given amount from the base stat.
func (s *Stat) SubBase(value int) {
	s.SetBase(s.base - value)
}

// GetValue returns the current value of the stat with mods applied.
func (s *Stat) GetValue() int {
	value := s.GetBase()
	for _, mod := range s.mods {
		value += mod
	}
	return min(s.GetDef().GetMax(), max(s.GetDef().GetMin(), value))
}

// Reset removes all mods.
func (s *Stat) Reset() {
	s.mods = make(map[int]int)
}

// SetMod applies a mod from a given source.
func (s *Stat) SetMod(source any, value int) error {
	statDef := s.GetDef()
	sourceObj, err := statDef.ctx.GetObject(source)
	if err != nil {
		return err
	}

	combatant, err := s.def.ctx.GetCombatantWithStat(s)
	if err == nil && combatant != nil {
		evt := &event.CombatantStatMod{CombatantID: combatant.GetID(), StatDefID: statDef.GetID(), SourceID: sourceObj.GetID(), ModValue: value}
		statDef.ctx.EmitEvent(evt)
		value = evt.ModValue
	}

	if s.mods == nil {
		s.mods = make(map[int]int)
	}
	s.mods[sourceObj.GetID()] = value
	if value == 0 {
		delete(s.mods, sourceObj.GetID())
	}
	return nil
}
