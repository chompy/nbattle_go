package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
)

type TestEffect struct {
	hpStat *nbattle.StatDef
}

func (e *TestEffect) OnAdd(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
	stat, _ := effectCtx.Target.GetStat(e.hpStat)
	stat.SetBase(29)
}

func (e *TestEffect) OnRemove(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
	stat, _ := effectCtx.Target.GetStat(e.hpStat)
	stat.SetBase(31)
}

func (e *TestEffect) OnEvent(ctx *nbattle.Context, effectCtx *nbattle.EffectContext, evt event.Event) {
	if evt.Type() == event.TickEvent {
		stat, _ := effectCtx.Target.GetStat(e.hpStat)
		stat.AddBase(-1)
	}
}

func TestCombatantEffect(t *testing.T) {

	ctx := nbattle.New()

	statDefHP := ctx.NewStatDef("hp", 0, 99)
	effectDefTest := ctx.NewEffectDef("test", func() nbattle.Effect {
		return &TestEffect{hpStat: statDefHP}
	})

	cmbt := ctx.NewCombatant()
	hpStat, err := cmbt.GetStat(statDefHP)
	if err != nil {
		t.Fatal(err)
	}
	hpStat.SetBase(15)

	srcCmbt := ctx.NewCombatant()

	cmbt.SetEffect(effectDefTest, srcCmbt, 1)

	hpStat, err = cmbt.GetStat(statDefHP)
	if err != nil {
		t.Fatal(err)
	}
	if hpStat.GetValue() != 29 {
		t.Fatal("expected effect to set combatant hp to 29 on add")
	}

	ctx.Tick()

	hpStat, err = cmbt.GetStat(statDefHP)
	if err != nil {
		t.Fatal(err)
	}
	if hpStat.GetValue() != 28 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	ctx.Tick()

	hpStat, err = cmbt.GetStat(statDefHP)
	if err != nil {
		t.Fatal(err)
	}
	if hpStat.GetValue() != 27 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	cmbt.SetEffect(effectDefTest, srcCmbt, 0)

	hpStat, err = cmbt.GetStat(statDefHP)
	if err != nil {
		t.Fatal(err)
	}

	ctx.Tick()
	if cmbt.HasEffect(hpStat) {
		t.Fatal("expected effect to be removed")
	}

	if hpStat.GetValue() != 31 {
		t.Fatal("expected effect to set combatant hp to 31 on removal")
	}

}
