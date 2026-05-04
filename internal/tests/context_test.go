package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go"
)

func initCtx() *nbattle.Context {
	ctx := nbattle.New()
	ctx.NewStatDef("max_hp", 0, 99)
	ctx.NewStatDef("hp", 0, 99)
	return ctx
}

func TestHandleEvent(t *testing.T) {

	srcCtx := initCtx()
	recCtx := initCtx()

	srcCtx.HookEvents(recCtx.HandleEvent)

	cmbt := srcCtx.NewCombatant()
	cmbt.GetStat("hp").SetBase(30)
	cmbt.GetStat("max_hp").SetBase(30)

	recCmbt, _ := recCtx.GetCombatantByID(cmbt.GetID())

	if cmbt.GetID() != recCmbt.GetID() {
		t.Fatal("receiving combatant id does not match source")
	}

	if cmbt.GetStat("hp").GetValue() != recCmbt.GetStat("hp").GetValue() {
		t.Fatal("receiving combatant hp does not match source")
	}

	cmbt.GetStat("hp").AddBase(-5)
	if cmbt.GetStat("hp").GetValue() != recCmbt.GetStat("hp").GetValue() {
		t.Fatal("receiving combatant hp does not match source")
	}

}
