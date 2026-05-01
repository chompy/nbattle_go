package event

type RemoveCombatantEffect struct {
	TargetID    int
	EffectDefID int
}

func (e *RemoveCombatantEffect) Type() Type { return RemoveCombatantEffectEvent }

func (e *RemoveCombatantEffect) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.TargetID, e.EffectDefID)
}

func (e *RemoveCombatantEffect) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	targetId, err := d.ReadInt()
	if err != nil {
		return err
	}
	effectDefId, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.TargetID = targetId
	e.EffectDefID = effectDefId
	return nil
}
