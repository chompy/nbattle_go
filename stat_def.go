package nbattle

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
