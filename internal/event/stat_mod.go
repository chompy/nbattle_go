package event

type StatMod struct {
	StatID   int
	SourceID int
	ModValue int
}

func (e *StatMod) Type() Type { return StatModEvent }

func (e *StatMod) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.StatID, e.SourceID, e.ModValue)
}

func (e *StatMod) Deserialize(data []byte) error {
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
	srcId, err := d.ReadInt()
	if err != nil {
		return err
	}
	val, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.StatID = statId
	e.SourceID = srcId
	e.ModValue = val
	return nil
}
