package nbattle

import "github.com/chompy/nbattle_go/internal/event"

type Stat struct {
	BaseObject
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
	s.base = min(s.GetDef().GetMax(), max(s.GetDef().GetMin(), value))
	s.ctx.EmitEvent(&event.StatBase{StatID: s.GetID(), Value: value})
}

func (s *Stat) AddBase(value int) {
	s.SetBase(s.base + value)
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
	sourceObj := s.ctx.GetObject(source)
	if sourceObj == nil {
		return
	}
	if s.mods == nil {
		s.mods = make(map[int]int)
	}
	s.mods[sourceObj.GetID()] = value
	if value == 0 {
		delete(s.mods, sourceObj.GetID())
	}
	s.ctx.EmitEvent(&event.StatMod{StatID: s.GetID(), SourceID: sourceObj.GetID(), ModValue: value})
}
