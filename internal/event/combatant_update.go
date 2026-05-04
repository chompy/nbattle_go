package event

type CombatantUpdate struct {
	CombatantID int
	Active      bool
}

func (e *CombatantUpdate) Type() Type { return CombatantUpdateEvent }

func (e *CombatantUpdate) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.CombatantID, e.Active)
}

func (e *CombatantUpdate) Deserialize(data []byte) error {
	d := newDeserializer(data)
	evType, err := d.ReadByte()
	if err != nil {
		return err
	}
	if Type(evType) != e.Type() {
		return ErrDeserializeWrongType
	}
	combatantID, err := d.ReadInt()
	if err != nil {
		return err
	}
	active, err := d.ReadBool()
	if err != nil {
		return err
	}

	e.CombatantID = combatantID
	e.Active = active
	return nil
}
