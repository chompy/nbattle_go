package nbattle

type StatDef struct {
	objectBase
	min int
	max int
}

func (d *StatDef) Type() objectType {
	return objectTypeStatDef
}

func (d *StatDef) Serialize() []byte {
	return nil
}

func (d *StatDef) Min() int {
	return d.min
}

func (d *StatDef) Max() int {
	return d.max
}
