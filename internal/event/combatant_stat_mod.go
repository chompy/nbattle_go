package event

type CombatantStatMod struct {
	CombatantID int
	StatDefID   int
	SourceID    int
	ModValue    int
}

func (e *CombatantStatMod) Type() Type { return CombatantStatModEvent }

func (e *CombatantStatMod) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.CombatantID, e.StatDefID, e.SourceID, e.ModValue)
}

func (e *CombatantStatMod) Deserialize(data []byte) error {
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
	statDefID, err := d.ReadInt()
	if err != nil {
		return err
	}
	sourceID, err := d.ReadInt()
	if err != nil {
		return err
	}

	modValue, err := d.ReadInt()
	if err != nil {
		return err
	}
	e.CombatantID = combatantID
	e.StatDefID = statDefID
	e.SourceID = sourceID
	e.ModValue = modValue
	return nil
}
