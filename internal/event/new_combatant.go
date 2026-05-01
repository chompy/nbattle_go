package event

type NewCombatant struct {
	ID   int
	Team int
}

func (e *NewCombatant) Type() Type {
	return NewCombatantEvent
}

func (e *NewCombatant) Serialize() ([]byte, error) {
	return serialize(NewCombatantEvent, e.ID, e.Team)
}

func (e *NewCombatant) Deserialize(data []byte) error {
	d := newDeserializer(data)
	eventType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(eventType) != NewCombatantEvent {
		return ErrDeserializeWrongType
	}
	combatantID, err := d.ReadInt()
	if err != nil {
		return err
	}
	team, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.ID = combatantID
	e.Team = team
	return nil
}
