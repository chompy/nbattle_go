package event

type CombatantFlag struct {
	TargetID int
	Flag     uint64
	On       bool
}

func (e *CombatantFlag) Type() Type { return CombatantFlagEvent }

func (e *CombatantFlag) Serialize() ([]byte, error) {
	return serialize(e.Type(), e.TargetID, e.Flag, e.On)
}

func (e *CombatantFlag) Deserialize(data []byte) error {
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
	flag, err := d.ReadUint64()
	if err != nil {
		return err
	}
	on, err := d.ReadBool()
	if err != nil {
		return err
	}
	e.TargetID = targetID
	e.Flag = flag
	e.On = on
	return nil
}
