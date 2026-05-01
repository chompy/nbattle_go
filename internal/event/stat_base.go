package event

type StatBase struct {
	StatID int
	Value  int
}

func (e *StatBase) Type() Type { return StatBaseEvent }

func (e *StatBase) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.StatID, e.Value)
}

func (e *StatBase) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	statId, err := d.ReadInt()
	if err != nil {
		return err
	}
	val, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.StatID = statId
	e.Value = val
	return nil
}
