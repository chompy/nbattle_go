package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go/nbattle_go"
)

type TestEffect struct {
	source *nbattle.Combatant
	target *nbattle.Combatant
}

func (e *TestEffect) Source() *nbattle.Combatant {
	return e.source
}

func (e *TestEffect) Target() *nbattle.Combatant {
	return e.target
}

func (e *TestEffect) OnApply() {

}

func (e *TestEffect) OnRemove() {

}

func (e *TestEffect) OnEvent(event *nbattle.Event) {

}

func TestCombatantEffect(t *testing.T) {
}
