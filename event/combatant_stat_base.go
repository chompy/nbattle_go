package event

type CombatantStatBase struct {
	CombatantID int
	StatDefID   int
	Value       int
}

func (e *CombatantStatBase) Type() Type { return CombatantStatBaseEvent }

func (e *CombatantStatBase) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.CombatantID, e.StatDefID, e.Value)
}

func (e *CombatantStatBase) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	combatantId, err := d.ReadInt()
	if err != nil {
		return err
	}
	statDefId, err := d.ReadInt()
	if err != nil {
		return err
	}
	value, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.CombatantID = combatantId
	e.StatDefID = statDefId
	e.Value = value
	return nil
}
