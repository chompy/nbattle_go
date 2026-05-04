package nbattle

import "github.com/chompy/nbattle_go/internal/event"

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

func (d *StatDef) GetMin() int {
	return d.min
}

func (d *StatDef) GetMax() int {
	return d.max
}

type Stat struct {
	def  *StatDef
	base int
	mods map[int]int
}

func (s *Stat) GetDef() *StatDef {
	return s.def
}

func (s *Stat) GetBase() int {
	return s.base
}

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

func (s *Stat) AddBase(value int) {
	s.SetBase(s.base + value)
}

func (s *Stat) SubBase(value int) {
	s.SetBase(s.base - value)
}

func (s *Stat) GetValue() int {
	value := s.GetBase()
	for _, mod := range s.mods {
		value += mod
	}
	return min(s.GetDef().GetMax(), max(s.GetDef().GetMin(), value))
}

func (s *Stat) Reset() {
	s.mods = make(map[int]int)
}

func (s *Stat) SetMod(source any, value int) {
	statDef := s.GetDef()

	sourceObj := statDef.ctx.GetObject(source)
	if sourceObj == nil {
		return
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
}
