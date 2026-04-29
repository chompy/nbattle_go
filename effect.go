package nbattle

type EffectDef struct {
	objectBase
	create func(target *Combatant, source *Combatant) Effect
}

func (d *EffectDef) Type() objectType {
	return objectTypeEffectDef
}

func (d *EffectDef) Serialize() []byte {
	return nil
}

type Effect interface {
	Source() *Combatant
	Target() *Combatant
	OnApply()
	OnRemove()
	OnEvent(event *Event)
}
