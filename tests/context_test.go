package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go/nbattle_go"
)

func TestHandleEvent(t *testing.T) {

	srcCtx := nbattle.New()
	recCtx := nbattle.New()

	srcCtx.HookEvents(recCtx.HandleEvent)

	statDefHP := srcCtx.NewStatDef(0, 99)
	statDefMaxHP := srcCtx.NewStatDef(0, 99)
	recCtx.NewStatDef(0, 99)
	recCtx.NewStatDef(0, 99)

	cmbt := srcCtx.NewCombatant()
	cmbt.Stat(statDefHP).SetBase(30)
	cmbt.Stat(statDefMaxHP).SetBase(30)

	recCmbt, _ := recCtx.GetCombatantByID(cmbt.ID())

	if cmbt.ID() != recCmbt.ID() {
		t.Fatal("receiving combatant id does not match source")
	}

	if cmbt.Stat(statDefHP).Value() != recCmbt.Stat(statDefHP).Value() {
		t.Fatal("receiving combatant hp does not match source")
	}

	cmbt.Stat(statDefHP).AddBase(-5)
	if cmbt.Stat(statDefHP).Value() != recCmbt.Stat(statDefHP).Value() {
		t.Fatal("receiving combatant hp does not match source")
	}

}
