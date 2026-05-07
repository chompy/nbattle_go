package event

type Trigger struct {
	TriggerDefID   int
	EffectDefID    int
	EffectTargetID int
	EffectSourceID int
	EffectPotency  int
}

func (e *Trigger) Type() Type {
	return TriggerEvent
}

func (e *Trigger) Serialize() ([]byte, error) {
	return serialize(TriggerEvent, e.TriggerDefID, e.EffectDefID, e.EffectTargetID, e.EffectSourceID, e.EffectPotency)
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
	effectDefID, err := d.ReadInt()
	if err != nil {
		return err
	}
	effectTargetID, err := d.ReadInt()
	if err != nil {
		return err
	}
	effectSourceID, err := d.ReadInt()
	if err != nil {
		return err
	}
	effectPotency, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.TriggerDefID = triggerDefID
	e.EffectDefID = effectDefID
	e.EffectTargetID = effectTargetID
	e.EffectSourceID = effectSourceID
	e.EffectPotency = effectPotency
	return nil
}
