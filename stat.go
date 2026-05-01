package nbattle

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
	//s.ctx.EmitEvent(EventTypeStatBase, s.GetID(), value)
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

func (s *Stat) Mod(source any, value int) {
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
	//s.ctx.EmitEvent(EventTypeStatMod, s.GetID(), sourceObj.GetID(), value)
}
