package event

type Hook func(e Event) error

type Type uint8

const (
	TickEvent Type = iota
	CombatantUpdateEvent
	CombatantStatBaseEvent
	CombatantStatModEvent
	CombatantEffectEvent
)

type Event interface {
	Type() Type
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}
