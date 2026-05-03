package nbattle

import "github.com/chompy/nbattle_go/internal/event"

type Stat struct {
	BaseObject
	def  *StatDef
	base int
	mods map[int]int
}

func (s *Stat) GetType() ObjectType {
	return ObjectTypeStat
}

func (s *Stat) GetDef() *StatDef {
	return s.def
}

func (s *Stat) GetBase() int {
	return s.base
}

func (s *Stat) SetBase(value int) {
	evt := &event.StatBase{StatID: s.GetID(), Value: value}
	s.ctx.EmitEvent(evt)
	s.base = min(s.GetDef().GetMax(), max(s.GetDef().GetMin(), evt.Value))
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
	evt := &event.StatMod{StatID: s.GetID(), SourceID: sourceObj.GetID(), ModValue: value}
	s.ctx.EmitEvent(evt)
	if s.mods == nil {
		s.mods = make(map[int]int)
	}
	s.mods[sourceObj.GetID()] = evt.ModValue
	if evt.ModValue == 0 {
		delete(s.mods, sourceObj.GetID())
	}
}
