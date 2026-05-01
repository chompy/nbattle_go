package event

type Tick struct {
	Tick int
}

func (e *Tick) Type() Type {
	return TickEvent
}

func (e *Tick) Serialize() ([]byte, error) {
	return serialize(TickEvent, e.Tick)
}

func (e *Tick) Deserialize(data []byte) error {
	d := newDeserializer(data)
	eventType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(eventType) != TickEvent {
		return ErrDeserializeWrongType
	}
	tick, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.Tick = tick
	return nil
}
