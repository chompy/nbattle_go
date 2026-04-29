package nbattle

type Stat struct {
	objectBase
	def  *StatDef
	base int
	mods map[int]int
}

func (d *Stat) Type() objectType {
	return objectTypeStat
}

func (d *Stat) Serialize() []byte {
	return nil
}

func (s *Stat) Def() *StatDef {
	return s.def
}

func (s *Stat) Base() int {
	return s.base
}

func (s *Stat) SetBase(value int) {
	s.base = min(s.Def().Max(), max(s.Def().Min(), value))
	s.ctx.EmitEvent(EventTypeStatBase, s.ID(), value)
}

func (s *Stat) AddBase(value int) {
	s.SetBase(s.base + value)
}

func (s *Stat) Value() int {
	value := s.Base()
	for _, mod := range s.mods {
		value += mod
	}
	return min(s.Def().Max(), max(s.Def().Min(), value))
}

func (s *Stat) Reset() {
	s.mods = make(map[int]int)
}

func (s *Stat) Mod(source object, value int) {
	if s.mods == nil {
		s.mods = make(map[int]int)
	}
	s.mods[source.ID()] = value
	if value == 0 {
		delete(s.mods, source.ID())
	}
	s.ctx.EmitEvent(EventTypeStatMod, s.ID(), source.ID(), value)
}
