package event

type CombatantEffect struct {
	TargetID    int
	EffectDefID int
	Potency     int
	SourceID    int
}

func (e *CombatantEffect) Type() Type { return CombatantEffectEvent }

func (e *CombatantEffect) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.TargetID, e.EffectDefID, e.Potency, e.SourceID)
}

func (e *CombatantEffect) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	targetID, err := d.ReadInt()
	if err != nil {
		return err
	}
	effectDefID, err := d.ReadInt()
	if err != nil {
		return err
	}
	potency, err := d.ReadInt()
	if err != nil {
		return err
	}
	sourceID, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.TargetID = targetID
	e.EffectDefID = effectDefID
	e.Potency = potency
	e.SourceID = sourceID
	return nil
}
