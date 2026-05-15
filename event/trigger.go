package event

type Trigger struct {
	TriggerDefID int
	SourceID     int
}

func (e *Trigger) Type() Type {
	return TriggerEvent
}

func (e *Trigger) Serialize() ([]byte, error) {
	return serialize(TriggerEvent, e.TriggerDefID, e.SourceID)
}

func (e *Trigger) Deserialize(data []byte) error {
	d := newDeserializer(data)
	eventType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(eventType) != TriggerEvent {
		return ErrDeserializeWrongType
	}
	triggerDefID, err := d.ReadInt()
	if err != nil {
		return err
	}
	sourceID, err := d.ReadInt()
	if err != nil {
		return err
	}

	e.TriggerDefID = triggerDefID
	e.SourceID = sourceID
	return nil
}
