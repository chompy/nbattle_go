package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go/internal/combat"
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
	hpStat, err := cmbt.GetStat("hp")
	if err != nil {
		t.Fatal(err)
	}
	hpStat.SetBase(30)
	maxHpStat, err := cmbt.GetStat("max_hp")
	if err != nil {
		t.Fatal(err)
	}
	maxHpStat.SetBase(30)

	recCmbt, err := recCtx.GetCombatantByID(cmbt.GetID())
	if err != nil {
		t.Fatal(err)
	}

	if cmbt.GetID() != recCmbt.GetID() {
		t.Fatal("receiving combatant id does not match source")
	}

	srcHp, _ := cmbt.GetStat("hp")
	recHp, _ := recCmbt.GetStat("hp")

	if srcHp.GetValue() != recHp.GetValue() {
		t.Fatal("receiving combatant hp does not match source")
	}

	hpStat, _ = cmbt.GetStat("hp")
	hpStat.AddBase(-5)
	srcHp, _ = cmbt.GetStat("hp")
	recHp, _ = recCmbt.GetStat("hp")
	if srcHp.GetValue() != recHp.GetValue() {
		t.Fatal("receiving combatant hp does not match source")
	}

}
