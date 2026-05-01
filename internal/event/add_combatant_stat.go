package event

type AddCombatantStat struct {
	CombatantID int
	StatID      int
	StatDefID   int
}

func (e *AddCombatantStat) Type() Type { return AddCombatantStatEvent }

func (e *AddCombatantStat) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.CombatantID, e.StatID, e.StatDefID)
}

func (e *AddCombatantStat) Deserialize(data []byte) error {
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
	statId, err := d.ReadInt()
	if err != nil {
		return err
	}
	statDefId, err := d.ReadInt()
	if err != nil {
		return err
	}

	e.CombatantID = combatantId
	e.StatID = statId
	e.StatDefID = statDefId
	return nil
}
