package event

type AddCombatantEffect struct {
	EffectDefID int
	TargetID    int
	SourceID    int
}

func (e *AddCombatantEffect) Type() Type { return AddCombatantEffectEvent }

func (e *AddCombatantEffect) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.EffectDefID, e.TargetID, e.SourceID)
}

func (e *AddCombatantEffect) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	effectDefId, err := d.ReadInt()
	if err != nil {
		return err
	}
	targetId, err := d.ReadInt()
	if err != nil {
		return err
	}
	sourceId, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.EffectDefID = effectDefId
	e.TargetID = targetId
	e.SourceID = sourceId
	return nil
}
