package nbattle

type EventHook func(e *Event) error

type EventType uint8

const (
	EventTypeTick                  EventType = iota // tick(int)
	EventTypeStatBase                               // stat(int),base(int)
	EventTypeStatMod                                // stat(int),source(int),mod(int)
	EventTypeCombatantNew                           // combatant(int)
	EventTypeCombatantStatAdd                       // combatant(int),stat(int),statDef(int)
	EventTypeCombatantEffectApply                   // target(int),effectDef(int),source(int)
	EventTypeCombatantEffectRemove                  // target(int),effectDef(int)
)

type Event struct {
	eventType EventType
	tick      int
	values    []any
}

func (e *Event) Type() EventType {
	return e.eventType
}

func (e *Event) Get(index int) any {
	if index < 0 || index > len(e.values)-1 {
		return nil
	}
	return e.values[index]
}

func (e *Event) GetInt(index int) int {
	switch v := e.Get(index).(type) {
	case int:
		return v
	}
	return 0
}
