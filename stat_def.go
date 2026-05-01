package nbattle

type StatDef struct {
	BaseObject
	name string
	min  int
	max  int
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
